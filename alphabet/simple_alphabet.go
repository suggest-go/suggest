package alphabet

// simpleAlphabet is alphabet, which use map for mapping char to int
type simpleAlphabet struct {
	set map[rune]struct{}
	charHolder
}

// NewSimpleAlphabet returns new instance of SimpleAlphabet
func NewSimpleAlphabet(chars []rune) Alphabet {
	set := make(map[rune]struct{}, len(chars))

	for _, char := range chars {
		set[char] = struct{}{}
	}

	return &simpleAlphabet{
		set:        set,
		charHolder: charHolder{chars},
	}
}

//
func (a *simpleAlphabet) Has(char rune) bool {
	_, ok := a.set[char]
	return ok
}
