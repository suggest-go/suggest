package suggest

import (
	"fmt"
	"sync"

	"github.com/suggest-go/suggest/pkg/analysis"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/utils"

	"github.com/suggest-go/suggest/pkg/index"
	"golang.org/x/sync/errgroup"
)

// Suggester is the interface that provides the access to
// approximate string search
type Suggester interface {
	// Suggest returns top-k similar candidates
	Suggest(query string, similarity float64, metric metric.Metric, factory CollectorManagerFactory) ([]Candidate, error)
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
func (n *nGramSuggester) Suggest(query string, similarity float64, metric metric.Metric, factory CollectorManagerFactory) ([]Candidate, error) {
	tokens := n.tokenizer.Tokenize(query)

	if len(tokens) == 0 {
		return []Candidate{}, nil
	}

	sizeA := len(tokens)
	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	lenIndices := n.indices.Size()

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	// channel that receives fuzzyCollector and performs a search on length segment
	sizeCh := make(chan int, bMax-bMin+1)
	workerPool := errgroup.Group{}
	collectorManager := factory()

	locker := sync.Mutex{}
	similarityHolder := utils.AtomicFloat64{}
	similarityHolder.Store(similarity)

	for i := 0; i < utils.Min(maxSearchQueriesAtOnce, bMax-bMin+1); i++ {
		workerPool.Go(func() error {
			for sizeB := range sizeCh {
				threshold := metric.Threshold(similarityHolder.Load(), sizeA, sizeB)

				// it means that the similarity has been changed and we will skip this value processing
				if threshold == 0 || threshold > sizeB || threshold > sizeA {
					continue
				}

				invertedIndex := n.indices.Get(sizeB)

				if invertedIndex == nil {
					continue
				}

				collector := collectorManager.Create()
				collector.SetScorer(NewMetricScorer(metric, sizeA, sizeB))

				if err := n.searcher.Search(invertedIndex, tokens, threshold, collector); err != nil {
					return fmt.Errorf("failed to search posting lists: %w", err)
				}

				locker.Lock()

				if err := collectorManager.Collect(collector); err != nil {
					locker.Unlock()

					return err
				}

				if fuzzy, ok := collectorManager.(*FuzzyCollectorManager); ok && fuzzy.GetLowestScore() > similarityHolder.Load() {
					similarityHolder.Store(fuzzy.GetLowestScore())
				}

				locker.Unlock()
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

	return collectorManager.GetCandidates(), nil
}
