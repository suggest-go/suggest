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

// trouble with mapping, we sort alphabet in size order for it, but will be collision when size are equals
func (a *compositeAlphabet) MapChar(char rune) int32 {
	key := int32(InvalidChar)
	shift := 0
	for _, alphabet := range a.alphabets {
		key = alphabet.MapChar(char)
		if key != InvalidChar {
			key += int32(shift)
			break
		}

		shift += alphabet.Size()
	}

	return key
}
