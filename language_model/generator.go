package language_model

type NGram = []Word
type NGrams = []NGram

type Generator interface {
	Generate(sentence Sentence) NGrams
}

func NewGenerator(
	nGramOrder uint8,
	startSymbol, endSymbol Word,
) *generator {
	return &generator{
		nGramOrder:  nGramOrder,
		startSymbol: startSymbol,
		endSymbol:   endSymbol,
	}
}

//
type generator struct {
	nGramOrder             uint8
	startSymbol, endSymbol Word
}

//
func (g *generator) Generate(sentence Sentence) NGrams {
	nGrams := NGrams{}
	k := int(g.nGramOrder)

	sentence = append([]Word{g.startSymbol}, sentence...)
	sentence = append(sentence, g.endSymbol)

	for i := 0; i <= len(sentence)-k; i++ {
		nGrams = append(nGrams, sentence[i:i+k])
	}

	return nGrams
}
