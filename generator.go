package suggest

const maxN = 8

type Generator interface {
	Generate(word string) []Term
	NGramSize() int
}

type generatorImpl struct {
	nGramSize int
	alphabet Alphabet
}

func NewGenerator(nGramSize int, alphabet Alphabet) Generator {
	return &generatorImpl{nGramSize, alphabet}
}

func (g *generatorImpl) Generate(word string) []Term {
	nGrams := splitIntoNGrams(word, g.nGramSize)
	l := len(nGrams)
	set := make(map[Term]struct{}, l)
	list := make([]Term, 0, l)
	for _, nGram := range nGrams {
		index := g.nGramToIndex(nGram)
		_, found := set[index]
		set[index] = struct{}{}
		if !found {
			list = append(list, index)
		}
	}

	return list
}

func (g *generatorImpl) NGramSize() int {
	return g.nGramSize
}

// SplitIntoNGrams split given word on k-gram list
// inspired by https://github.com/Lazin/go-ngram
func splitIntoNGrams(word string, k int) []string {
	sliceLen := len(word) - k + 1
	if sliceLen <= 0 || sliceLen > len(word) {
		panic("Invalid word length for spliting")
	}

	var prevIndexes [maxN]int
	result := make([]string, 0, sliceLen)
	i := 0
	for index := range word {
		i++
		if i > k {
			top := prevIndexes[(i-k)%k]
			substr := word[top:index]
			result = append(result, substr)
		}

		prevIndexes[i%k] = index
	}

	top := prevIndexes[(i+1)%k]
	substr := word[top:]
	result = append(result, substr)

	return result
}

// Map ngram to int (index)
func (g *generatorImpl) nGramToIndex(nGram string) Term {
	index := int32(0)
	size := g.alphabet.Size()
	for _, char := range nGram {
		i := g.alphabet.MapChar(char)
		if index == InvalidChar {
			panic("Invalid char was detected")
		}

		index = index*int32(size) + i
	}

	return Term(index)
}
