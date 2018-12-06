package alphabet

// russianAlphabet represents russian alphabet а-я
type russianAlphabet struct {
	parent Alphabet
}

// NewRussianAlphabet returns new instance of RussianAlphabet
func NewRussianAlphabet() Alphabet {
	return &russianAlphabet{
		NewSequentialAlphabet('а', 'я'),
	}
}

// Note, that we map ё as e
func (a *russianAlphabet) Has(char rune) bool {
	if char == 'ё' {
		return a.parent.Has('е')
	}

	return a.parent.Has(char)
}

func (a *russianAlphabet) Size() int {
	return a.parent.Size()
}

func (a *russianAlphabet) Chars() []rune {
	return a.parent.Chars()
}
