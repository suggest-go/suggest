package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/merger"
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
	merger    merger.ListMerger
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(
	cleaner index.Cleaner,
	generator index.Generator,
	indices index.InvertedIndexIndices,
	merger merger.ListMerger,
) NGramIndex {
	return &nGramIndexImpl{
		cleaner:   cleaner,
		indices:   indices,
		generator: generator,
		merger:    merger,
	}
}

// Suggest returns top-k similar candidates
func (n *nGramIndexImpl) Suggest(config *SearchConfig) ([]Candidate, error) {
	preparedQuery := n.cleaner.CleanAndWrap(config.query)

	return n.fuzzySearch(preparedQuery, config, NewTopKSelector(config.topK))
}

// AutoComplete returns candidates where the query string is a substring of each candidate
func (n *nGramIndexImpl) AutoComplete(query string, limit int) ([]Candidate, error) {
	preparedQuery := n.cleaner.CleanAndLeftWrap(query)

	return n.autoComplete(preparedQuery, NewTopKSelector(limit))
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

	rid := make([]index.PostingList, 0, len(set))
	sizeA := len(set)
	metric := config.metric
	similarity := config.similarity

	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	lenIndices := n.indices.Size()

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	boundarySlice := make([]int, 0, 2)
	i, j := sizeA, sizeA

	// iteratively expand the slice with boundaries [i, j] in the interval [bMin, bMax]
	for {
		boundarySlice = boundarySlice[:0]

		if i >= bMin {
			boundarySlice = append(boundarySlice, i)
		}

		if j != i && j <= bMax {
			boundarySlice = append(boundarySlice, j)
		}

		j++
		i--

		if len(boundarySlice) == 0 {
			break
		}

		for _, sizeB := range boundarySlice {
			threshold := metric.Threshold(similarity, sizeA, sizeB)

			// this is an unusual case, we should skip it
			if threshold == 0 {
				continue
			}

			// no reason to continue (the lowest candidate is more suitaible)
			if !selector.CanTakeWithScore(1 - metric.Distance(sizeA, sizeA, sizeB)) {
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

			// maximum allowable nGram miss count
			allowedSkips := sizeA - threshold + 1

			for _, term := range set {
				if allowedSkips == 0 {
					break
				}

				if !invertedIndex.Has(term) {
					allowedSkips--
				}
			}

			// no reason to continue, we have already reached all allowed skips
			if allowedSkips == 0 {
				continue
			}

			rid = rid[:0]

			// maybe run it concurrent?
			// go func() { buildRid, mergeCandidates, ch <- {candidates, sizeA, sizeB}
			// in main goroutine just collect it
			for _, term := range set {
				postingList, err := invertedIndex.Get(term)

				if err != nil {
					return nil, fmt.Errorf("Failed to retrieve a posting list: %v", err)
				}

				if len(postingList) > 0 {
					rid = append(rid, postingList)
				}
			}

			candidates := n.merger.Merge(rid, threshold)

			for _, c := range candidates {
				score := 1 - metric.Distance(c.Overlap, sizeA, sizeB)
				selector.Add(c.Position, score)
			}
		}
	}

	return selector.GetCandidates(), nil
}

// autoComplete performs a completion of phrases that contain the given query
func (n *nGramIndexImpl) autoComplete(query string, selector TopKSelector) ([]Candidate, error) {
	set := n.generator.Generate(query)
	rid := make([]index.PostingList, 0, len(set))
	invertedIndex := n.indices.GetWholeIndex()

	for _, term := range set {
		if !invertedIndex.Has(term) {
			return []Candidate{}, nil
		}
	}

	for _, term := range set {
		postingList, err := invertedIndex.Get(term)

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve a posting list: %v", err)
		}

		rid = append(rid, postingList)
	}

	for i, c := range n.merger.Merge(rid, len(rid)) {
		selector.Add(c.Position, float64(-i))
	}

	return selector.GetCandidates(), nil
}
