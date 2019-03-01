package language_model

const bufferSize = 100

type NGramBuilder struct {
	indexer                    Indexer
	startSymbolId, endSymbolId WordId
}

func NewNGramBuilder(
	indexer Indexer,
	startSymbol, endSymbol string,
) *NGramBuilder {
	return &NGramBuilder{
		indexer:       indexer,
		startSymbolId: indexer.GetOrCreate(startSymbol),
		endSymbolId:   indexer.GetOrCreate(endSymbol),
	}
}

func (nb *NGramBuilder) Build(retriever SentenceRetriver, nGramOrder uint8) CountingTrie {
	trie := NewTrie()
	seq := make([]WordId, 0, bufferSize)
	ch, quit := nb.produce(retriever)

	for {
		select {
		case sentence := <-ch:
			seq = seq[:0]
			seq = append(seq, nb.startSymbolId)

			for _, word := range sentence {
				seq = append(seq, nb.indexer.GetOrCreate(word))
			}

			seq = append(seq, nb.endSymbolId)

			for k := 1; k <= int(nGramOrder); k++ {
				for i := 0; i <= len(seq)-k; i++ {
					trie.Put(seq[i : i+k])
				}
			}

		case <-quit:
			return trie
		}
	}
}

func (nb *NGramBuilder) produce(retriever SentenceRetriver) (chan Sentence, chan bool) {
	ch := make(chan Sentence)
	quit := make(chan bool)

	go func() {
		for {
			sentence := retriever.Retrieve()
			if sentence == nil {
				break
			}

			if len(sentence) == 0 {
				continue
			}

			ch <- sentence
		}

		quit <- true
	}()

	return ch, quit
}
