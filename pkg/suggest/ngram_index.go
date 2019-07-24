package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/utils"

	"github.com/alldroll/suggest/pkg/index"
)

// NGramIndex is a structure that provides an access to
// approximate string search and autocomplete
type NGramIndex interface {
	// Suggest returns top-k similar candidates
	Suggest(config *SearchConfig) ([]Candidate, error)
	// AutoComplete returns candidates with query as substring
	AutoComplete(query string, limit int) ([]Candidate, error)
}

// nGramIndexImpl implements NGramIndex
type nGramIndexImpl struct {
	cleaner   index.Cleaner
	indices   index.InvertedIndexIndices
	generator index.Generator
	searcher  index.Searcher
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(
	cleaner index.Cleaner,
	generator index.Generator,
	indices index.InvertedIndexIndices,
	searcher index.Searcher,
) NGramIndex {
	return &nGramIndexImpl{
		cleaner:   cleaner,
		indices:   indices,
		generator: generator,
		searcher:  searcher,
	}
}

// Suggest returns top-k similar candidates
func (n *nGramIndexImpl) Suggest(config *SearchConfig) ([]Candidate, error) {
	preparedQuery := n.cleaner.CleanAndWrap(config.query)

	return n.fuzzySearch(preparedQuery, config, NewTopKSelectorWithRanker(config.topK, &idOrderRank{}))
}

// AutoComplete returns candidates where the query string is a substring of each candidate
func (n *nGramIndexImpl) AutoComplete(query string, limit int) ([]Candidate, error) {
	preparedQuery := n.cleaner.CleanAndLeftWrap(query)

	return n.autoComplete(preparedQuery, NewTopKSelectorWithRanker(limit, &idOrderRank{}))
}

// fuzzySearch performs approximate string search in the search index
func (n *nGramIndexImpl) fuzzySearch(
	query string,
	config *SearchConfig,
	selector TopKSelector,
) ([]Candidate, error) {
	set := n.generator.Generate(query)

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

// autoComplete performs a completion of phrases that contain the given query
func (n *nGramIndexImpl) autoComplete(query string, selector TopKSelector) ([]Candidate, error) {
	set := n.generator.Generate(query)
	invertedIndex := n.indices.GetWholeIndex()
	candidates, err := n.searcher.Search(invertedIndex, set, len(set))

	if err != nil {
		return nil, fmt.Errorf("failed to search posting lists: %v", err)
	}

	for i, c := range candidates {
		selector.Add(c.Position(), float64(-i))
	}

	return selector.GetCandidates(), nil
}
