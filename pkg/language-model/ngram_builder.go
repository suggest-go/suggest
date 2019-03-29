package lm

const bufferSize = 100

// NGramBuilder is an entity that responsible for creating CountTrie
type NGramBuilder struct {
	indexer                    Indexer
	startSymbolID, endSymbolID WordID
}

// NewNGramBuilder returns new instance of NGramBuilder
func NewNGramBuilder(
	indexer Indexer,
	startSymbol, endSymbol string,
) *NGramBuilder {
	return &NGramBuilder{
		indexer:       indexer,
		startSymbolID: indexer.GetOrCreate(startSymbol),
		endSymbolID:   indexer.GetOrCreate(endSymbol),
	}
}

// Build builds CountTrie with nGrams
func (nb *NGramBuilder) Build(retriever SentenceRetriever, nGramOrder uint8) CountTrie {
	trie := NewCountTrie()
	seq := make([]WordID, 0, bufferSize)
	ch, quit := nb.produce(retriever)

	for {
		select {
		case sentence := <-ch:
			seq = seq[:0]
			seq = append(seq, nb.startSymbolID)

			for _, word := range sentence {
				seq = append(seq, nb.indexer.GetOrCreate(word))
			}

			seq = append(seq, nb.endSymbolID)

			for k := 1; k <= int(nGramOrder); k++ {
				for i := 0; i <= len(seq)-k; i++ {
					trie.Put(seq[i:i+k], 1)
				}
			}

		case <-quit:
			return trie
		}
	}
}

// Transfers retrieved sentences to channel
func (nb *NGramBuilder) produce(retriever SentenceRetriever) (chan Sentence, chan bool) {
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
