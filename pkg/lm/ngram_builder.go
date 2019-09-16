package lm

// NGramBuilder is an entity that responsible for creating CountTrie
type NGramBuilder struct {
	startSymbol, endSymbol Token
}

// NewNGramBuilder returns new instance of NGramBuilder
func NewNGramBuilder(
	startSymbol, endSymbol string,
) *NGramBuilder {
	return &NGramBuilder{
		startSymbol: startSymbol,
		endSymbol:   endSymbol,
	}
}

// Build builds CountTrie with nGrams
func (nb *NGramBuilder) Build(retriever SentenceRetriever, nGramOrder uint8) CountTrie {
	trie := NewCountTrie()
	ch, quit := nb.produce(retriever)

	for {
		select {
		case sentence := <-ch:
			sentence = append([]Token{nb.startSymbol}, sentence...)
			sentence = append(sentence, nb.endSymbol)

			for k := 1; k <= int(nGramOrder); k++ {
				for i := 0; i <= len(sentence)-k; i++ {
					trie.Put(sentence[i:i+k], 1)
				}
			}

		case <-quit:
			return trie
		}
	}
}

// produce transfers retrieved sentences to the channel
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
