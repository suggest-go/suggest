package suggest

import "github.com/alldroll/suggest/alphabet"

type Alphabet = alphabet.Alphabet

// NewSimpleAlphabet returns new instance of SimpleAlphabet
func NewSimpleAlphabet(chars []rune) Alphabet {
	return alphabet.NewSimpleAlphabet(chars)
}

// NewSequentialAlphabet returns new instance of sequentialAlphabet
func NewSequentialAlphabet(min, max rune) Alphabet {
	return alphabet.NewSequentialAlphabet(min, max)
}

// NewNumberAlphabet returns new instance of numberAlphabet
func NewNumberAlphabet() Alphabet {
	return alphabet.NewNumberAlphabet()
}

// NewEnglishAlphabet returns new instance of englishAlphabet
func NewEnglishAlphabet() Alphabet {
	return alphabet.NewEnglishAlphabet()
}

// NewRussianAlphabet returns new instance of RussianAlphabet
func NewRussianAlphabet() Alphabet {
	return alphabet.NewRussianAlphabet()
}

// NewCompositeAlphabet returns new instance of compositeAlphabet
func NewCompositeAlphabet(alphabets []Alphabet) Alphabet {
	return alphabet.NewCompositeAlphabet(alphabets)
}
