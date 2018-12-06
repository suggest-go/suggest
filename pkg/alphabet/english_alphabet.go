package alphabet

// englishAlphabet represents english aphabet
type englishAlphabet struct {
	Alphabet
}

// NewEnglishAlphabet returns new instance of englishAlphabet
func NewEnglishAlphabet() Alphabet {
	return &englishAlphabet{
		NewSequentialAlphabet('a', 'z'),
	}
}
