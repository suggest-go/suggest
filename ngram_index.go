package suggest

import (
	"github.com/alldroll/suggest/index"
	"github.com/alldroll/suggest/list_merger"
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
	AutoComplete(query string, topK int) []Candidate
}

// nGramIndexImpl implements NGramIndex
type nGramIndexImpl struct {
	cleaner   index.Cleaner
	indices   index.InvertedIndexIndices
	generator index.Generator
	merger    list_merger.ListMerger
}

// NewNGramIndex returns a new NGramIndex object
func NewNGramIndex(cleaner index.Cleaner, generator index.Generator, indices index.InvertedIndexIndices, merger list_merger.ListMerger) NGramIndex {
	return &nGramIndexImpl{
		cleaner:   cleaner,
		indices:   indices,
		generator: generator,
		merger:    merger,
	}
}

// Suggest returns top-k similar strings
func (n *nGramIndexImpl) Suggest(config *SearchConfig) []FuzzyCandidate {
	result := make([]FuzzyCandidate, 0, config.topK)
	preparedQuery := n.cleaner.Clean(config.query)
	if len(preparedQuery) < 3 { // TODO дичь
		return result
	}

	return n.fuzzySearch(preparedQuery, config)
}

// AutoComplete returns candidates with query as substring
func (n *nGramIndexImpl) AutoComplete(query string, topK int) []Candidate {
	return nil
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
	collector := NewTopKCollector(topK)

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	type pp struct {
		candidates   []*list_merger.MergeCandidate
		sizeA, sizeB int
	}

	for sizeB := bMax; sizeB >= bMin; sizeB-- {
		threshold := metric.Threshold(similarity, sizeA, sizeB)
		if threshold == 0 {
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
			collector.Add(c.Pos, distance)
		}
	}

	return collector.GetCandidates()
}
