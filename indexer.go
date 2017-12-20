package suggest

type Index []map[Term]PostingList

type Indexer interface {
	Index(dictionary Dictionary) Index
}

func NewIndexer(
	nGramSize int,
	generator Generator,
	cleaner Cleaner,
) Indexer {
	return &indexerImpl{
		nGramSize,
		generator,
		cleaner,
	}
}

type indexerImpl struct {
	nGramSize int
	generator Generator
	cleaner Cleaner
}

func (ix *indexerImpl) Index(dictionary Dictionary) Index {
	i := dictionary.Iterator()
	indices := make(Index, 0)

	for {
		key, word := i.GetPair()

		if len(word) >= ix.nGramSize {

			prepared := ix.cleaner.Clean(word)
			set := ix.generator.Generate(prepared)
			cardinality := len(set)

			if len(indices) <= cardinality {
				tmp := make(Index, cardinality+1, cardinality*2)
				copy(tmp, indices)
				indices = tmp
			}

			table := indices[cardinality]
			if table == nil {
				table = make(map[Term]PostingList)
				indices[cardinality] = table
			}

			for _, index := range set {
				table[index] = append(table[index], key)
			}
		}

		if !i.Next() {
			break
		}
	}

	return indices
}
