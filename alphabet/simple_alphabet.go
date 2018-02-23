package alphabet

// simpleAlphabet is alphabet, which use map for mapping char to int
type simpleAlphabet struct {
	table map[rune]int32
	charHolder
}

// NewSimpleAlphabet returns new instance of SimpleAlphabet
func NewSimpleAlphabet(chars []rune) Alphabet {
	table := make(map[rune]int32, len(chars))
	index := int32(0)
	for _, char := range chars {
		table[char] = index
		index++
	}

	return &simpleAlphabet{
		table:      table,
		charHolder: charHolder{chars},
	}
}

func (a *simpleAlphabet) MapChar(char rune) int32 {
	index, ok := a.table[char]
	if !ok {
		index = InvalidChar
	}

	return index
}
