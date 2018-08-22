package alphabet

import "sort"

// compositeAlphabet represents composite pattern for a group of alphabets
type compositeAlphabet struct {
	alphabets []Alphabet
	charHolder
}

// NewCompositeAlphabet returns new instance of compositeAlphabet
func NewCompositeAlphabet(alphabets []Alphabet) Alphabet {
	size := 0

	sort.Slice(alphabets, func(i, j int) bool {
		return alphabets[i].Size() < alphabets[j].Size()
	})

	for _, alphabet := range alphabets {
		size += alphabet.Size()
	}

	chars := make([]rune, 0, size)
	for _, alphabet := range alphabets {
		chars = append(chars, alphabet.Chars()...)
	}

	return &compositeAlphabet{
		alphabets:  alphabets,
		charHolder: charHolder{chars},
	}
}

//
func (a *compositeAlphabet) Has(char rune) bool {
	has := false

	for _, alphabet := range a.alphabets {
		if has = alphabet.Has(char); has {
			break
		}
	}

	return has
}
