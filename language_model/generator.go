package language_model

type NGram = []Word
type NGrams = []NGram

type Generator interface {
	Generate() NGrams
}

func NewGenerator(
	sentenceRetriever SentenceRetriver,
	nGramOrder uint8,
	startSymbol, endSymbol Word,
) *generator {
	return &generator{
		sentenceRetriever: sentenceRetriever,
		nGramOrder:        nGramOrder,
		startSymbol:       startSymbol,
		endSymbol:         endSymbol,
	}
}

//
type generator struct {
	sentenceRetriever      SentenceRetriver
	nGramOrder             uint8
	startSymbol, endSymbol Word
}

//
func (g *generator) Generate() NGrams {
	nGrams := NGrams{}
	k := int(g.nGramOrder)
	sentence := g.sentenceRetriever.Retrieve()

	if sentence == nil {
		return nil
	}

	sentence = append([]Word{g.startSymbol}, sentence...)
	sentence = append(sentence, g.endSymbol)

	for i := 0; i <= len(sentence)-k; i++ {
		nGrams = append(nGrams, sentence[i:i+k])
	}

	return nGrams
}
