package suggest

import (
	"fmt"
	"sync/atomic"

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
	similarityVal := atomic.Value{}
	similarityVal.Store(config.similarity)

	topKQueue := NewTopKQueue(config.topK)
	// done notifies that
	done := make(chan bool)
	// channel that receives the finished collectors
	colOutCh := make(chan Collector)

	// update similarityVal value if the received score is greater than similarityVal
	// this optimization prevents scanning the segments, where we couldn't reach the same similarityVal
	go func() {
		for collector := range colOutCh {
			// fill the topKQueue with new candidates
			for _, item := range collector.GetCandidates() {
				topKQueue.Add(item.Key, item.Score)
			}

			// if we have achieved the top-k local candidates,
			// try to update the similarity value with the lowest canidates score
			if topKQueue.IsFull() {
				similarity := similarityVal.Load().(float64)
				score := topKQueue.GetLowestScore()

				if similarity < score {
					similarityVal.Store(score)
				}
			}
		}

		done <- true
	}()

	// channel that receives fuzzyCollector and performs a search on length segment
	colInCh := make(chan *fuzzyCollector)
	workerPool := errgroup.Group{}

	for i := 0; i < utils.Min(maxSearchQueriesAtOnce, bMax-bMin+1); i++ {
		workerPool.Go(func() error {
			for collector := range colInCh {
				similarity := similarityVal.Load().(float64)
				threshold := config.metric.Threshold(similarity, sizeA, collector.sizeB)

				// this is an unusual case, we should skip it
				if threshold == 0 || threshold > collector.sizeB {
					continue
				}

				invertedIndex := n.indices.Get(collector.sizeB)

				if invertedIndex == nil {
					continue
				}

				if err := n.searcher.Search(invertedIndex, set, threshold, collector); err != nil {
					return fmt.Errorf("failed to search posting lists: %v", err)
				}

				colOutCh <- collector
			}

			return nil
		})
	}

	buf := [2]int{}

	// iteratively expand the slice with boundaries [i, j] in the interval [bMin, bMax]
	for i, j := sizeA, sizeA; i >= bMin || j <= bMax; i, j = i-1, j+1 {
		boundarySlice := buf[:0]

		if i >= bMin {
			boundarySlice = append(boundarySlice, i)
		}

		if j != i && j <= bMax {
			boundarySlice = append(boundarySlice, j)
		}

		for _, sizeB := range boundarySlice {
			collector := &fuzzyCollector{
				metric:    config.metric,
				sizeA:     sizeA,
				sizeB:     sizeB,
				topKQueue: NewTopKQueue(config.topK),
			}

			colInCh <- collector
		}
	}

	// close collector in channel
	close(colInCh)

	if err := workerPool.Wait(); err != nil {
		close(colOutCh)
		return nil, err
	}

	// close collector out channel
	close(colOutCh)
	// wait for result collected
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
