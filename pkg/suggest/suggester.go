package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/analysis"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/utils"

	"github.com/alldroll/suggest/pkg/index"
	"golang.org/x/sync/errgroup"
)

// Suggester is the interface that provides the access to
// approximate string search
type Suggester interface {
	// Suggest returns top-k similar candidates
	Suggest(config *SearchConfig) ([]Candidate, error)
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
func (n *nGramSuggester) Suggest(config *SearchConfig) ([]Candidate, error) {
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

	topKQueue := NewTopKQueue(config.topK)
	// done notifies that all subqueries were processed and we have top-k candidates
	done := make(chan bool)
	// channel that receives the collected candidates
	subResultCh := make(chan []Candidate)

	// update similarityVal value if the received score is greater than similarityHolders value
	// this optimization prevents scanning the segments, where we couldn't reach the same similarityVal
	go func() {
		for subResult := range subResultCh {
			// fill the topKQueue with new candidates
			for _, item := range subResult {
				topKQueue.Add(item.Key, item.Score)
			}

			// if we have achieved the top-k local candidates,
			// try to update the similarity value with the lowest candidates score
			if topKQueue.IsFull() {
				score := topKQueue.GetLowestScore()

				if similarityHolder.Load() < score {
					similarityHolder.Store(score)
				}
			}
		}

		done <- true
	}()

	// channel that receives fuzzyCollector and performs a search on length segment
	sizeCh := make(chan int)
	workerPool := errgroup.Group{}

	for i := 0; i < utils.Min(maxSearchQueriesAtOnce, bMax-bMin+1); i++ {
		workerPool.Go(func() error {
			for sizeB := range sizeCh {
				similarity := similarityHolder.Load()
				threshold := config.metric.Threshold(similarity, sizeA, sizeB)

				// this is an unusual case, we should skip it
				if threshold == 0 || threshold > sizeB {
					continue
				}

				invertedIndex := n.indices.Get(sizeB)

				if invertedIndex == nil {
					continue
				}

				collector := &fuzzyCollector{
					sizeA:     sizeA,
					sizeB:     sizeB,
					metric:    config.metric,
					topKQueue: NewTopKQueue(config.topK),
				}

				if err := n.searcher.Search(invertedIndex, set, threshold, collector); err != nil {
					return fmt.Errorf("failed to search posting lists: %v", err)
				}

				subResultCh <- collector.GetCandidates()
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
		close(subResultCh)
		return nil, err
	}

	// close collector out channel
	close(subResultCh)
	// wait for the result collected
	<-done

	return topKQueue.GetCandidates(), nil
}

type fuzzyCollector struct {
	metric    metric.Metric
	sizeA     int
	sizeB     int
	topKQueue TopKQueue
}

// Collect collects the given merge candidate
// calculates the distance, and tries to add this document with it's score to the collector
func (c *fuzzyCollector) Collect(item merger.MergeCandidate) error {
	score := 1 - c.metric.Distance(item.Overlap(), c.sizeA, c.sizeB)
	c.topKQueue.Add(item.Position(), score)

	return nil
}

// GetCandidates returns `top k items`
func (c *fuzzyCollector) GetCandidates() []Candidate {
	return c.topKQueue.GetCandidates()
}

// Score returns the score of the given position
func (c *fuzzyCollector) SetScorer(scorer Scorer) {
	return
}
