package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/analysis"
	"github.com/alldroll/suggest/pkg/utils"

	"github.com/alldroll/suggest/pkg/index"
)

// Suggester is the interface that provides the access to
// approximate string search
type Suggester interface {
	// Suggest returns top-k similar candidates
	Suggest(config *SearchConfig) ([]Candidate, error)
}

// nGramSuggester implements Suggester
type nGramSuggester struct {
	indices   index.InvertedIndexIndices
	searcher  index.Searcher
	tokenizer analysis.Tokenizer
	ranker    Rank
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
		ranker:    &idOrderRank{},
	}
}

// Suggest returns top-k similar candidates
func (n *nGramSuggester) Suggest(config *SearchConfig) ([]Candidate, error) {
	selector := NewTopKCollectorWithRanker(config.topK, n.ranker)
	set := n.tokenizer.Tokenize(config.query)

	if len(set) == 0 {
		return []Candidate{}, nil
	}

	sizeA := len(set)
	metric := config.metric
	similarity := config.similarity

	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	lenIndices := n.indices.Size()

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	buf := [2]int{}
	i, j := sizeA, sizeA

	// iteratively expand the slice with boundaries [i, j] in the interval [bMin, bMax]
	for {
		boundarySlice := buf[:0]

		if i >= bMin {
			boundarySlice = append(boundarySlice, i)
		}

		if j != i && j <= bMax {
			boundarySlice = append(boundarySlice, j)
		}

		if len(boundarySlice) == 0 {
			break
		}

		j++
		i--

		for _, sizeB := range boundarySlice {
			threshold := metric.Threshold(similarity, sizeA, sizeB)

			// this is an unusual case, we should skip it
			if threshold == 0 {
				continue
			}

			// no reason to continue (the lowest candidate is more suitable even if we have complete intersection)
			if !selector.CanTakeWithScore(1 - metric.Distance(utils.Min(sizeA, sizeB), sizeA, sizeB)) {
				continue
			}

			// maybe the lowest candidate's score in the selector will give us a bigger threshold
			// here we should use new threshold, only if selector has collected topK elements
			if selector.IsFull() {
				lowestScore := selector.GetLowestScore()
				thresholdForLowestScore := metric.Threshold(lowestScore, sizeA, sizeB)

				if thresholdForLowestScore > threshold {
					threshold = thresholdForLowestScore
				}
			}

			// no reason to continue: threshold is more than the right border of out viewable interval
			if threshold > sizeB {
				continue
			}

			invertedIndex := n.indices.Get(sizeB)

			if invertedIndex == nil {
				continue
			}

			candidates, err := n.searcher.Search(invertedIndex, set, threshold)

			if err != nil {
				return nil, fmt.Errorf("failed to search posting lists: %v", err)
			}

			for _, c := range candidates {
				score := 1 - metric.Distance(c.Overlap(), sizeA, sizeB)
				selector.Add(c.Position(), score)
			}
		}
	}

	return selector.GetCandidates(), nil
}
