package suggest

// InvalidChar represents unmapped value of given char
const InvalidChar = -1

// Alphabet is abstract for manipulating with set of symbols
type Alphabet interface {
	// MapChar map given char to int
	MapChar(char rune) int
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

// SimpleAlphabet is alphabet, which use map for mapping char to int
type SimpleAlphabet struct {
	table map[rune]int
	charHolder
}

// NewSimpleAlphabet returns new instance of SimpleAlphabet
func NewSimpleAlphabet(chars []rune) *SimpleAlphabet {
	table := make(map[rune]int, len(chars))
	index := 0
	for _, char := range chars {
		table[char] = index
		index++
	}

	return &SimpleAlphabet{table, charHolder{chars}}
}

// MapChar map given char to int
func (a *SimpleAlphabet) MapChar(char rune) int {
	index, ok := a.table[char]
	if !ok {
		index = InvalidChar
	}

	return index
}

// SequentialAlphabet represents alphabet with continuous list of ascii characters
type SequentialAlphabet struct {
	min, max rune
	charHolder
}

// NewSequentialAlphabet returns new instance of SequentialAlphabet
func NewSequentialAlphabet(min, max rune) *SequentialAlphabet {
	chars := make([]rune, 0, max-min+1)
	for ch := min; ch <= max; ch++ {
		chars = append(chars, ch)
	}

	return &SequentialAlphabet{
		min, max, charHolder{chars},
	}
}

// MapChar map given char to int
func (a *SequentialAlphabet) MapChar(char rune) int {
	if char < a.min || char > a.max {
		return InvalidChar
	}

	return int(char - a.min)
}

// RussianAlphabet represents russian alphabet а-я
type RussianAlphabet struct {
	parent *SequentialAlphabet
}

// NewRussianAlphabet returns new instance of RussianAlphabet
func NewRussianAlphabet() *RussianAlphabet {
	return &RussianAlphabet{
		NewSequentialAlphabet('а', 'я'),
	}
}

// MapChar map given char to int. Note, that we map ё as e
func (a *RussianAlphabet) MapChar(char rune) int {
	if char == 'ё' {
		return a.parent.MapChar('е')
	}

	return a.parent.MapChar(char)
}

// Size returns the size of alphabet
func (a *RussianAlphabet) Size() int {
	return a.parent.Size()
}

// Chars returns the current set of symbols
func (a *RussianAlphabet) Chars() []rune {
	return a.parent.Chars()
}

// EnglishAlphabet represents english aphabet
type EnglishAlphabet struct {
	*SequentialAlphabet
}

// NewEnglishAlphabet returns new instance of EnglishAlphabet
func NewEnglishAlphabet() *EnglishAlphabet {
	return &EnglishAlphabet{
		NewSequentialAlphabet('a', 'z'),
	}
}

// NumberAlphabet represents number alphabet [0-9]
type NumberAlphabet struct {
	*SequentialAlphabet
}

// NewNumberAlphabet returns new instance of NumberAlphabet
func NewNumberAlphabet() *NumberAlphabet {
	return &NumberAlphabet{
		NewSequentialAlphabet('0', '9'),
	}
}

// CompositeAlphabet represents composite pattern for a group of alphabets
type CompositeAlphabet struct {
	alphabets []Alphabet
	charHolder
}

// NewCompositeAlphabet returns new instance of CompositeAlphabet
func NewCompositeAlphabet(alphabets []Alphabet) *CompositeAlphabet {
	size := 0
	for _, alphabet := range alphabets {
		size += alphabet.Size()
	}

	chars := make([]rune, 0, size)
	for _, alphabet := range alphabets {
		chars = append(chars, alphabet.Chars()...)
	}

	return &CompositeAlphabet{alphabets, charHolder{chars}}
}

// MapChar map given char to int
func (a *CompositeAlphabet) MapChar(char rune) int {
	key := InvalidChar
	shift := 0
	for _, alphabet := range a.alphabets {
		key = alphabet.MapChar(char)
		if key != InvalidChar {
			key += shift
			break
		}

		shift += alphabet.Size()
	}

	return key
}
