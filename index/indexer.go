package index

import "github.com/alldroll/suggest/dictionary"

type Index = map[Term]PostingList

type Indices = []Index

type Indexer interface {
	IndexIndices(dictionary dictionary.Dictionary) Indices
	Index(dictionary dictionary.Dictionary) Index
}

func NewIndexer(
	nGramSize int,
	generator Generator,
	cleaner Cleaner,
) Indexer {
	return &indexerImpl{
		nGramSize: nGramSize,
		generator: generator,
		cleaner:   cleaner,
	}
}

type indexerImpl struct {
	nGramSize int
	generator Generator
	cleaner   Cleaner
}

func (ix *indexerImpl) IndexIndices(dictionary dictionary.Dictionary) Indices {
	i := dictionary.Iterator()
	indices := make(Indices, 0)

	for {
		key, word := i.GetPair()

		if len(word) >= ix.nGramSize {
			prepared := ix.cleaner.Clean(word)
			set := ix.generator.Generate(prepared)
			cardinality := len(set)

			if len(indices) <= cardinality {
				tmp := make(Indices, cardinality+1, cardinality*2)
				copy(tmp, indices)
				indices = tmp
			}

			index := indices[cardinality]
			if index == nil {
				index = make(Index)
				indices[cardinality] = index
			}

			for _, term := range set {
				index[term] = append(index[term], key)
			}
		}

		if !i.Next() {
			break
		}
	}

	return indices
}

func (ix *indexerImpl) Index(dictionary dictionary.Dictionary) Index {
	i := dictionary.Iterator()
	index := make(Index, 0)

	for {
		key, word := i.GetPair()

		if len(word) >= ix.nGramSize {
			prepared := ix.cleaner.Clean(word)
			set := ix.generator.Generate(prepared)

			for _, term := range set {
				index[term] = append(index[term], key)
			}
		}

		if !i.Next() {
			break
		}
	}

	return index
}
