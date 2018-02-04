package suggest

/*
 * inspired by
 *
 * http://www.chokkan.org/software/simstring/
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 * http://www.aclweb.org/anthology/C10-1096
 */

import (
	"container/heap"
)

// NGramIndex is structure ... describe me please
type NGramIndex interface {
	// Suggest returns top-k similar candidates
	Suggest(config *SearchConfig) []Candidate
	// AutoComplete returns candidates with query as substring
	AutoComplete(query string, topK int) []Candidate
}

// nGramIndexImpl implements NGramIndex
type nGramIndexImpl struct {
	cleaner   Cleaner
	indices   InvertedIndexIndices
	generator Generator
	merger ListMerger
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(cleaner Cleaner, generator Generator, indices InvertedIndexIndices, merger ListMerger, ) NGramIndex {
	return &nGramIndexImpl{
		cleaner: cleaner,
		indices: indices,
		generator: generator,
		merger: merger,
	}
}

// Suggest returns top-k similar strings
func (n *nGramIndexImpl) Suggest(config *SearchConfig) []Candidate {
	result := make([]Candidate, 0, config.topK)
	preparedQuery := n.cleaner.Clean(config.query)
	if len(preparedQuery) < 3 { // TODO дичь
		return result
	}

	candidates := n.fuzzySearch(preparedQuery, config)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append(
			[]Candidate{{r.pos, r.distance}},
			result...,
		)
	}

	return result
}

// AutoComplete returns candidates with query as substring
func (n *nGramIndexImpl) AutoComplete(query string, topK int) []Candidate {
	return nil
}

// fuzzySearch
func (n *nGramIndexImpl) fuzzySearch(query string, config *SearchConfig) *heapImpl {
	set := n.generator.Generate(query)
	sizeA := len(set)

	metric := config.metric
	similarity := config.similarity
	topK := config.topK

	h := newHeap(topK)
	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	rid := make([]PostingList, 0, sizeA)
	lenIndices := n.indices.Size()
	var r *rank

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	for sizeB := bMax; sizeB >= bMin; sizeB-- {
		threshold := metric.Threshold(similarity, sizeA, sizeB)
		if threshold == 0 {
			continue
		}

		// reset slice
		rid = rid[:0]
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

		for _, term := range set {
			postingList := invertedIndex.Get(term)
			if len(postingList) > 0 {
				rid = append(rid, postingList)
			}
		}

		candidates := n.merger.Merge(rid, threshold)
		// use heap search for finding top k items in a list efficiently
		// see http://stevehanov.ca/blog/index.php?id=122
		for _, c := range candidates {
			distance := metric.Distance(c.Overlap, sizeA, sizeB)

			if h.Len() < topK || h.Top().(*rank).distance > distance {
				if h.Len() == topK {
					r = heap.Pop(h).(*rank)
				} else {
					r = &rank{
						pos: 0,
						distance: 0.0,
					}
				}

				r.pos = c.Pos
				r.distance = distance
				heap.Push(h, r)
			}
		}
	}

	return h
}
