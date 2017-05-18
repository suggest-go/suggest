package suggest

//
const InvalidChar = -1

//
type Alphabet interface {
	MapChar(char rune) int
	Size() int
	Chars() []rune
}

//
type SimpleAlphabet struct {
	table map[rune]int
	chars []rune
}

//
func NewSimpleAlphabet(chars []rune) *SimpleAlphabet {
	table := make(map[rune]int, len(chars))
	index := 0
	for _, char := range chars {
		table[char] = index
		index++
	}

	return &SimpleAlphabet{table, chars}
}

//
func (a *SimpleAlphabet) MapChar(char rune) int {
	index, ok := a.table[char]
	if !ok {
		index = InvalidChar
	}

	return index
}

//
func (a *SimpleAlphabet) Chars() []rune {
	return a.chars
}

//
func (a *SimpleAlphabet) Size() int {
	return len(a.chars)
}

//
type SequentialAlphabet struct {
	chars    []rune
	min, max rune
}

//
func NewSequentialAlphabet(min, max rune) *SequentialAlphabet {
	chars := make([]rune, 0, max-min+1)
	for ch := min; ch <= max; ch++ {
		chars = append(chars, ch)
	}

	return &SequentialAlphabet{
		chars, min, max,
	}
}

//
func (a *SequentialAlphabet) MapChar(char rune) int {
	if char < a.min || char > a.max {
		return InvalidChar
	}

	return int(char - a.min)
}

//
func (a *SequentialAlphabet) Size() int {
	return len(a.chars)
}

//
func (a *SequentialAlphabet) Chars() []rune {
	return a.chars
}

//
type RussianAlphabet struct {
	parent *SequentialAlphabet
}

//
func NewRussianAlphabet() *RussianAlphabet {
	return &RussianAlphabet{
		NewSequentialAlphabet('а', 'я'),
	}
}

//
func (a *RussianAlphabet) MapChar(char rune) int {
	if char == 'ё' {
		return a.parent.MapChar('е')
	}

	return a.parent.MapChar(char)
}

//
func (a *RussianAlphabet) Size() int {
	return a.parent.Size()
}

//
func (a *RussianAlphabet) Chars() []rune {
	return a.parent.Chars()
}

//
type EnglishAlphabet struct {
	*SequentialAlphabet
}

func NewEnglishAlphabet() *EnglishAlphabet {
	return &EnglishAlphabet{
		NewSequentialAlphabet('a', 'z'),
	}
}

//
type NumberAlphabet struct {
	*SequentialAlphabet
}

//
func NewNumberAlphabet() *NumberAlphabet {
	return &NumberAlphabet{
		NewSequentialAlphabet('0', '9'),
	}
}

//
type CompositeAlphabet struct {
	alphabets []Alphabet
	chars     []rune
}

//
func NewCompositeAlphabet(alphabets []Alphabet) *CompositeAlphabet {
	size := 0
	for _, alphabet := range alphabets {
		size += alphabet.Size()
	}

	chars := make([]rune, 0, size)
	for _, alphabet := range alphabets {
		chars = append(chars, alphabet.Chars()...)
	}

	return &CompositeAlphabet{alphabets, chars}
}

//
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

//
func (a *CompositeAlphabet) Size() int {
	return len(a.chars)
}

//
func (a *CompositeAlphabet) Chars() []rune {
	return a.chars
}
