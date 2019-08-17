package index

const maxN = 8

// Generator represents entity for creating terms from given word
type Generator interface {
	// Generate returns terms array for given word
	Generate(word string) []Term
}

// generatorImpl implements Generator interface
type generatorImpl struct {
	nGramSize int
}

// NewGenerator returns new instance of Generator
func NewGenerator(nGramSize int) Generator {
	return &generatorImpl{
		nGramSize: nGramSize,
	}
}

// Generate returns terms array for given word
// inspired by https://github.com/Lazin/go-ngram
func (g *generatorImpl) Generate(word string) []Term {
	if len(word) < g.nGramSize {
		return []Term{}
	}

	result := make([]Term, 0, len(word)-g.nGramSize+1)
	prevIndexes := [maxN]int{}
	i := 0

	for index := range word {
		i++

		if i > g.nGramSize {
			top := prevIndexes[(i-g.nGramSize)%g.nGramSize]
			nGram := word[top:index]
			result = appendUnique(result, nGram)
		}

		prevIndexes[i%g.nGramSize] = index
	}

	top := prevIndexes[(i+1)%g.nGramSize]
	nGram := word[top:]
	result = appendUnique(result, nGram)

	return result
}

// https://blog.golang.org/profiling-go-programs
func appendUnique(a []Term, x Term) []Term {
	for _, y := range a {
		if x == y {
			return a
		}
	}

	return append(a, x)
}
