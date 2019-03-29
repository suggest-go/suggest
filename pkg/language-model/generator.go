package lm

type NGram = []Token
type NGrams = []NGram

// Generator is entity that responsible for transferring the given
// sentence into a sequence of NGrams
type Generator interface {
	// Generates a sequence of NGrams for the given sentence
	Generate(sentence Sentence) NGrams
}

// NewGenerator creates a generator entity
func NewGenerator(
	nGramOrder uint8,
	startSymbol, endSymbol Token,
) *generator {
	return &generator{
		nGramOrder:  nGramOrder,
		startSymbol: startSymbol,
		endSymbol:   endSymbol,
	}
}

// generator implements Generator interface
type generator struct {
	nGramOrder             uint8
	startSymbol, endSymbol Token
}

// Generates a sequence of NGrams for the given sentence
func (g *generator) Generate(sentence Sentence) NGrams {
	nGrams := NGrams{}
	k := int(g.nGramOrder)

	sentence = append([]Token{g.startSymbol}, sentence...)
	sentence = append(sentence, g.endSymbol)

	for i := 0; i <= len(sentence)-k; i++ {
		nGrams = append(nGrams, sentence[i:i+k])
	}

	return nGrams
}
