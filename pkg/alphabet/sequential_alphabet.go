package alphabet

// sequentialAlphabet represents alphabet with continuous list of ascii characters
type sequentialAlphabet struct {
	min, max rune
	charHolder
}

// NewSequentialAlphabet returns new instance of sequentialAlphabet
func NewSequentialAlphabet(min, max rune) Alphabet {
	chars := make([]rune, 0, max-min+1)
	for ch := min; ch <= max; ch++ {
		chars = append(chars, ch)
	}

	return &sequentialAlphabet{
		min:        min,
		max:        max,
		charHolder: charHolder{chars},
	}
}

func (a *sequentialAlphabet) Has(char rune) bool {
	return char >= a.min && char <= a.max
}
