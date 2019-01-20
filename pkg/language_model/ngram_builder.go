package language_model

const bufferSize = 100

type NGramBuilder struct {
	sentenceRetriever SentenceRetriver
	generators        []Generator
	indexer           Indexer
	nGramOrder        int
}

func NewNGramBuilder(
	sentenceRetriever SentenceRetriver,
	indexer Indexer,
	nGramOrder uint8,
) *NGramBuilder {
	return &NGramBuilder{
		sentenceRetriever: sentenceRetriever,
		indexer:           indexer,
		nGramOrder:        int(nGramOrder),
	}
}

func (nb *NGramBuilder) Build() Trie {
	trie := NewTrie()
	ch := make(chan Sentence)
	quit := make(chan bool)

	go func() {
		for {
			sentence := nb.sentenceRetriever.Retrieve()
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

	seq := make([]WordId, 0, bufferSize)

	for {
		select {
		case sentence := <-ch:
			seq = seq[:0]
			seq = append(seq, nb.indexer.GetOrCreate("<S>"))

			for _, word := range sentence {
				seq = append(seq, nb.indexer.GetOrCreate(word))
			}

			seq = append(seq, nb.indexer.GetOrCreate("<S/>"))

			for k := 1; k <= nb.nGramOrder; k++ {
				for i := 0; i <= len(seq)-k; i++ {
					trie.Put(seq[i : i+k])
				}
			}

		case <-quit:
			return trie
		}
	}
}
