package alphabet

// numberAlphabet represents number alphabet [0-9]
type numberAlphabet struct {
	Alphabet
}

// NewNumberAlphabet returns new instance of numberAlphabet
func NewNumberAlphabet() Alphabet {
	return &numberAlphabet{
		NewSequentialAlphabet('0', '9'),
	}
}
