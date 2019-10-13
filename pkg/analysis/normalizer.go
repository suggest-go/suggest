package analysis

import (
	"github.com/suggest-go/suggest/pkg/alphabet"
)

type normalizeFilter struct {
	chars alphabet.Alphabet
	pad   string
}

// NewNormalizerFilter returns tokens filter
func NewNormalizerFilter(chars alphabet.Alphabet, pad string) TokenFilter {
	return &normalizeFilter{
		chars: chars,
		pad:   pad,
	}
}

// Filter filters the given list with described behaviour
func (f *normalizeFilter) Filter(list []Token) []Token {
	for i, token := range list {
		res := ""

		for _, r := range token {
			if f.chars.Has(r) {
				res += string(r)
			} else {
				res += f.pad
			}
		}

		list[i] = res
	}

	return list
}
