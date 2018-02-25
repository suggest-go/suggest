package index

import "github.com/alldroll/suggest/dictionary"

type Index = map[Term]PostingList

type Indices = []Index

type Indexer interface {
	Index(dictionary dictionary.Dictionary) Indices
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

func (ix *indexerImpl) Index(dictionary dictionary.Dictionary) Indices {
	i := dictionary.Iterator()
	indices := make(Indices, 1)
	indices[0] = make(Index)

	for {
		key, word := i.GetPair()

		if len(word) >= ix.nGramSize {
			prepared := ix.cleaner.CleanAndWrap(word)
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
				indices[0][term] = append(indices[0][term], key)
			}
		}

		if !i.Next() {
			break
		}
	}

	return indices
}
