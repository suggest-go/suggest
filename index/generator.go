package index

import "github.com/alldroll/suggest/alphabet"

const maxN = 8

// Generator represents entity for creating terms from given word
type Generator interface {
	// Generate returns terms array for given word
	Generate(word string) []Term
}

// generatorImpl implements Generator interface
type generatorImpl struct {
	nGramSize int
	alphabet  alphabet.Alphabet
}

// NewGenerator returns new instance of Generator
func NewGenerator(nGramSize int, alphabet alphabet.Alphabet) Generator {
	return &generatorImpl{
		nGramSize: nGramSize,
		alphabet:  alphabet,
	}
}

// Generate returns terms array for given word
func (g *generatorImpl) Generate(word string) []Term {
	nGrams := SplitIntoNGrams(word, g.nGramSize)
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

// SplitIntoNGrams split given word on k-gram list
func SplitIntoNGrams(word string, k int) []string {
	runes := []rune(word)
	result := make([]string, 0, len(runes))

	for i := 0; i < len(runes)-k+1; i++ {
		result = append(result, string(runes[i:i+k]))
	}

	return result
}

// Map ngram to int (index)
func (g *generatorImpl) nGramToIndex(nGram string) Term {
	index := int32(0)
	size := g.alphabet.Size()
	for _, char := range nGram {
		i := g.alphabet.MapChar(char)
		if index == alphabet.InvalidChar {
			panic("Invalid char was detected")
		}

		index = index*int32(size) + i
	}

	return Term(index)
}
