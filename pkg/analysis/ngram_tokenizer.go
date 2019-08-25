package analysis

const maxN = 8

// NewNGramTokenizer creates a new instance of Tokenizer
func NewNGramTokenizer(nGramSize int) Tokenizer {
	return &nGramTokenizer{
		nGramSize: nGramSize,
	}
}

type nGramTokenizer struct {
	nGramSize int
}

// Tokenize splits the given text on a sequence of tokens
func (t *nGramTokenizer) Tokenize(text string) []Token {
	if len(text) < t.nGramSize {
		return []Token{}
	}

	result := make([]Token, 0, len(text)-t.nGramSize+1)
	prevIndexes := [maxN]int{}
	i := 0

	for index := range text {
		i++

		if i > t.nGramSize {
			top := prevIndexes[(i-t.nGramSize)%t.nGramSize]
			nGram := text[top:index]
			result = appendUnique(result, nGram)
		}

		prevIndexes[i%t.nGramSize] = index
	}

	top := prevIndexes[(i+1)%t.nGramSize]
	nGram := text[top:]
	result = appendUnique(result, nGram)

	return result
}

// https://blog.golang.org/profiling-go-programs
func appendUnique(a []Token, x Token) []Token {
	for _, y := range a {
		if x == y {
			return a
		}
	}

	return append(a, x)
}
