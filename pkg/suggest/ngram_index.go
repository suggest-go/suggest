package suggest

import (
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/list_merger"
)

/*
 * inspired by
 *
 * http://www.chokkan.org/software/simstring/
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 * http://www.aclweb.org/anthology/C10-1096
 */

// NGramIndex is structure ... describe me please
type NGramIndex interface {
	// Suggest returns top-k similar candidates
	Suggest(config *SearchConfig) []FuzzyCandidate
	// AutoComplete returns candidates with query as substring
	AutoComplete(query string, limit int) []Candidate
}

// nGramIndexImpl implements NGramIndex
type nGramIndexImpl struct {
	cleaner   index.Cleaner
	indices   index.InvertedIndexIndices
	generator index.Generator
	merger    list_merger.ListMerger
	intersect list_merger.ListIntersect
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(
	cleaner index.Cleaner,
	generator index.Generator,
	indices index.InvertedIndexIndices,
	merger list_merger.ListMerger,
	intersect list_merger.ListIntersect,
) NGramIndex {
	return &nGramIndexImpl{
		cleaner:   cleaner,
		indices:   indices,
		generator: generator,
		merger:    merger,
		intersect: intersect,
	}
}

// Suggest returns top-k similar strings
func (n *nGramIndexImpl) Suggest(config *SearchConfig) []FuzzyCandidate {
	result := make([]FuzzyCandidate, 0)
	preparedQuery := n.cleaner.CleanAndWrap(config.query)
	if len(preparedQuery) < 3 {
		return result
	}

	return n.fuzzySearch(preparedQuery, config)
}

// AutoComplete returns candidates with query as substring
func (n *nGramIndexImpl) AutoComplete(query string, limit int) []Candidate {
	result := make([]Candidate, 0)
	preparedQuery := n.cleaner.CleanAndLeftWrap(query)

	if len(preparedQuery) < 3 {
		return result
	}

	return n.autoComplete(preparedQuery, limit)
}

// fuzzySearch
func (n *nGramIndexImpl) fuzzySearch(query string, config *SearchConfig) []FuzzyCandidate {
	set := n.generator.Generate(query)
	rid := make([]index.PostingList, 0, len(set))
	sizeA := len(set)

	metric := config.metric
	similarity := config.similarity
	topK := config.topK

	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	lenIndices := n.indices.Size()
	selector := NewTopKSelector(topK)

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	boundaryValues := make([]int, 0, 2)
	i, j := sizeA, sizeA

	for {
		boundaryValues = boundaryValues[:0]

		if i >= bMin {
			boundaryValues = append(boundaryValues, i)
		}

		if j != i && j <= bMax {
			boundaryValues = append(boundaryValues, j)
		}

		j++
		i--

		if len(boundaryValues) == 0 {
			break
		}

		for _, sizeB := range boundaryValues {
			threshold := metric.Threshold(similarity, sizeA, sizeB)
			if threshold == 0 {
				continue
			}

			lowestCandidate, lowestDistance := selector.GetLowestRecord()
			if lowestCandidate != nil && selector.Size() == topK {
				// there is no reason
				if metric.Distance(sizeA, sizeA, sizeB) > lowestDistance {
					continue
				}

				thresholdForLowestDistance := metric.Threshold(1.0-lowestDistance, sizeA, sizeB)
				if thresholdForLowestDistance > threshold {
					threshold = thresholdForLowestDistance
				}
			}

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
				// there is no reason to continue, because of threshold
				if allowedSkips == 0 {
					break
				}

				if !invertedIndex.Has(term) {
					allowedSkips--
				}
			}

			if allowedSkips == 0 {
				continue
			}

			rid = rid[:0]

			// maybe run it concurrent?
			// go func() { buildRid, mergeCandidates, ch <- {candidates, sizeA, sizeB}
			// in main goroutine just collect it
			for _, term := range set {
				postingList := invertedIndex.Get(term)
				if len(postingList) > 0 {
					rid = append(rid, postingList)
				}
			}

			candidates := n.merger.Merge(rid, threshold)

			for _, c := range candidates {
				distance := metric.Distance(c.Overlap, sizeA, sizeB)
				selector.Add(c, distance)
			}
		}
	}

	return selector.GetCandidates()
}

// autoComplete
func (n *nGramIndexImpl) autoComplete(query string, limit int) []Candidate {
	set := n.generator.Generate(query)
	rid := make([]index.PostingList, 0, len(set))
	result := make([]Candidate, 0)

	invertedIndex := n.indices.GetWholeIndex()
	if invertedIndex == nil {
		return result
	}

	for _, term := range set {
		if !invertedIndex.Has(term) {
			return result
		}
	}

	for _, term := range set {
		postingList := invertedIndex.Get(term)
		rid = append(rid, postingList)
	}

	for _, c := range n.intersect.Intersect(rid, limit) {
		result = append(result, Candidate{c.Position})
	}

	return result
}
