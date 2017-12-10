package suggest

type Index []map[int]PostingList

type Indexer interface {
	Index() Index
}

func NewIndexer(
	nGramSize int,
	dictionary Dictionary,
	generator Generator,
	cleaner Cleaner,
) Indexer {
	return &indexerImpl{
		nGramSize,
		dictionary,
		generator,
		cleaner,
	}
}

type indexerImpl struct {
	nGramSize int
	dictionary Dictionary
	generator Generator
	cleaner Cleaner
}

func (ix *indexerImpl) Index() Index {
	i := ix.dictionary.Iterator()
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
				table = make(map[int]PostingList)
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

/*
func (b *invertedIndexIndicesBuilderCDBImpl) store(writer cdb.Writer, table map[int]InvertedList) error {
	k := make([]byte, 4)

	for key, list := range table {
		if list == nil {
			continue
		}

		value := make([]byte, len(list) * 4)
		for i, x := range list {
			binary.LittleEndian.PutUint32(value[4*i:], uint32(x))
		}

		binary.LittleEndian.PutUint32(k, uint32(key))

		writer.Put(k, value)
	}

	return writer.Close()
}
*/
