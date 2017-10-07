package suggest

/*
 * inspired by
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 * http://www.aclweb.org/anthology/C10-1096
 */

import (
	"container/heap"
)

// NGramIndex is structure ... describe me please
type NGramIndex struct {
	clean      Cleaner
	indices    InvertedListsIndices
	generator  Generator
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(cleaner Cleaner, generator Generator, indices InvertedListsIndices) *NGramIndex {
	return &NGramIndex{
		cleaner, indices, generator,
	}
}

// Suggest returns top-k similar strings
func (n *NGramIndex) Suggest(config *SearchConfig) []Candidate {
	result := make([]Candidate, 0, config.topK)
	preparedQuery := n.clean.Clean(config.query)
	if len(preparedQuery) < 3 {
		return result
	}

	candidates := n.search(preparedQuery, config)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append(
			[]Candidate{{n.indices.IndexToWordKey(r.id), r.distance}},
			result...,
		)
	}

	return result
}

func (n *NGramIndex) search(query string, config *SearchConfig) *heapImpl {
	set := n.generator.Generate(query)
	sizeA := len(set)

	metric := config.metric
	similarity := config.similarity
	topK := config.topK

	h := &heapImpl{}
	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	rid := make([][]int, 0, sizeA)
	lenIndices := n.indices.Size()

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
		invertedLists := n.indices.Get(sizeB)
		// maximum allowable ngram miss count
		allowedSkips := sizeA - threshold + 1
		for _, index := range set {
			// there is no reason to continue, because of threshold
			if allowedSkips == 0 {
				break
			}

			list := invertedLists[index]
			if len(list) > 0 {
				rid = append(rid, list)
			} else {
				allowedSkips--
			}
		}

		if len(rid) < threshold {
			continue
		}

		counts := n.getCounts(rid, threshold)
		// use heap search for finding top k items in a list efficiently
		// see http://stevehanov.ca/blog/index.php?id=122
		var r *rank
		for inter := len(counts) - 1; inter >= threshold; inter-- {
			for _, id := range counts[inter] {
				distance := metric.Distance(inter, sizeA, sizeB)

				if h.Len() < topK || h.Top().(*rank).distance > distance {
					if h.Len() == topK {
						r = heap.Pop(h).(*rank)
					} else {
						r = &rank{0, 0}
					}

					r.id = id
					r.distance = distance
					heap.Push(h, r)
				}
			}
		}
	}

	return h
}

// TODO describe me!
func (n *NGramIndex) getCounts(rid [][]int, threshold int) [][]int {
	if threshold == 1 {
		return scanCount(rid, threshold)
	}

	return mergeSkip(rid, threshold)
}
