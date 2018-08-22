package alphabet

// Alphabet is abstract for manipulating with set of symbols
type Alphabet interface {
	// Has tells is given char exists in alphabet
	Has(char rune) bool
	// Size returns the size of alphabet
	Size() int
	// Chars returns the current set of symbols
	Chars() []rune
}

type charHolder struct {
	chars []rune
}

// Size returns the size of alphabet
func (c *charHolder) Size() int {
	return len(c.chars)
}

// Chars returns the current set of symbols
func (c *charHolder) Chars() []rune {
	return c.chars
}
