package suggest

import (
	"fmt"
	"sync"

	"github.com/suggest-go/suggest/pkg/analysis"
	"github.com/suggest-go/suggest/pkg/utils"

	"github.com/suggest-go/suggest/pkg/index"
	"golang.org/x/sync/errgroup"
)

// Suggester is the interface that provides the access to
// approximate string search
type Suggester interface {
	// Suggest returns top-k similar candidates
	Suggest(config SearchConfig) ([]Candidate, error)
}

// maxSearchQueriesAtOnce tells how many goroutines can be used at once for a search query
const maxSearchQueriesAtOnce = 5

// nGramSuggester implements Suggester
type nGramSuggester struct {
	indices   index.InvertedIndexIndices
	searcher  index.Searcher
	tokenizer analysis.Tokenizer
}

// NewSuggester returns a new Suggester instance
func NewSuggester(
	indices index.InvertedIndexIndices,
	searcher index.Searcher,
	tokenizer analysis.Tokenizer,
) Suggester {
	return &nGramSuggester{
		indices:   indices,
		searcher:  searcher,
		tokenizer: tokenizer,
	}
}

// Suggest returns top-k similar candidates
func (n *nGramSuggester) Suggest(config SearchConfig) ([]Candidate, error) {
	set := n.tokenizer.Tokenize(config.query)

	if len(set) == 0 {
		return []Candidate{}, nil
	}

	sizeA := len(set)
	bMin, bMax := config.metric.MinY(config.similarity, sizeA), config.metric.MaxY(config.similarity, sizeA)
	lenIndices := n.indices.Size()

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	// store similarity as atomic value
	// we are going to update its value after search sub-work complete
	similarityHolder := utils.AtomicFloat64{}
	similarityHolder.Store(config.similarity)

	topKQueue := topKQueuePool.Get().(TopKQueue)
	topKQueue.Reset(config.topK)
	defer topKQueuePool.Put(topKQueue)

	// channel that receives fuzzyCollector and performs a search on length segment
	sizeCh := make(chan int, bMax-bMin+1)
	workerPool := errgroup.Group{}
	lock := sync.Mutex{}

	for i := 0; i < utils.Min(maxSearchQueriesAtOnce, bMax-bMin+1); i++ {
		workerPool.Go(func() error {
			for sizeB := range sizeCh {
				similarity := similarityHolder.Load()
				threshold := config.metric.Threshold(similarity, sizeA, sizeB)

				// it means that the similarity has been changed and we will skip this value processing
				if threshold == 0 || threshold > sizeB || threshold > sizeA {
					continue
				}

				invertedIndex := n.indices.Get(sizeB)

				if invertedIndex == nil {
					continue
				}

				queue := topKQueuePool.Get().(TopKQueue)
				queue.Reset(config.topK)

				collector := &fuzzyCollector{
					topKQueue: queue,
					scorer:    NewMetricScorer(config.metric, sizeA, sizeB),
				}

				if err := n.searcher.Search(invertedIndex, set, threshold, collector); err != nil {
					return fmt.Errorf("failed to search posting lists: %v", err)
				}

				lock.Lock()

				topKQueue.Merge(queue)

				if topKQueue.IsFull() && similarityHolder.Load() < topKQueue.GetLowestScore() {
					similarityHolder.Store(topKQueue.GetLowestScore())
				}

				lock.Unlock()

				topKQueuePool.Put(queue)
			}

			return nil
		})
	}

	// iteratively expand the slice with boundaries [i, j] in the interval [bMin, bMax]
	for i, j := sizeA, sizeA+1; i >= bMin || j <= bMax; i, j = i-1, j+1 {
		if i >= bMin {
			sizeCh <- i
		}

		if j <= bMax {
			sizeCh <- j
		}
	}

	// close input channel for worker pool
	close(sizeCh)

	if err := workerPool.Wait(); err != nil {
		return nil, err
	}

	return topKQueue.GetCandidates(), nil
}

var topKQueuePool = sync.Pool{
	New: func() interface{} {
		return NewTopKQueue(50)
	},
}
