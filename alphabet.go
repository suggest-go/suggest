package suggest

const INVALID_CHAR = -1

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

func NewSimpleAlphabet(chars []rune) *SimpleAlphabet {
	table := make(map[rune]int, len(chars))
	index := 0
	for _, char := range chars {
		table[char] = index
		index++
	}

	return &SimpleAlphabet{table, chars}
}

func (self *SimpleAlphabet) MapChar(char rune) int {
	index, ok := self.table[char]
	if !ok {
		index = INVALID_CHAR
	}

	return index
}

func (self *SimpleAlphabet) Chars() []rune {
	return self.chars
}

func (self *SimpleAlphabet) Size() int {
	return len(self.chars)
}

//
type SequentialAlphabet struct {
	chars    []rune
	min, max rune
}

func NewSequentialAlphabet(min, max rune) *SequentialAlphabet {
	chars := make([]rune, 0, max-min+1)
	for ch := min; ch <= max; ch++ {
		chars = append(chars, ch)
	}

	return &SequentialAlphabet{
		chars, min, max,
	}
}

func (self *SequentialAlphabet) MapChar(char rune) int {
	if char < self.min || char > self.max {
		return INVALID_CHAR
	}

	return int(char - self.min)
}

func (self *SequentialAlphabet) Size() int {
	return len(self.chars)
}

func (self *SequentialAlphabet) Chars() []rune {
	return self.chars
}

type RussianAlphabet struct {
	parent *SequentialAlphabet
}

//
func NewRussianAlphabet() *RussianAlphabet {
	return &RussianAlphabet{
		NewSequentialAlphabet('а', 'я'),
	}
}

func (self *RussianAlphabet) MapChar(char rune) int {
	if char == 'ё' {
		return self.parent.MapChar('е')
	}

	return self.parent.MapChar(char)
}

func (self *RussianAlphabet) Size() int {
	return self.parent.Size()
}

func (self *RussianAlphabet) Chars() []rune {
	return self.parent.Chars()
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

func (self *CompositeAlphabet) MapChar(char rune) int {
	key := INVALID_CHAR
	shift := 0
	for _, alphabet := range self.alphabets {
		key = alphabet.MapChar(char)
		if key != INVALID_CHAR {
			key += shift
			break
		}

		shift += alphabet.Size()
	}

	return key
}

func (self *CompositeAlphabet) Size() int {
	return len(self.chars)
}

func (self *CompositeAlphabet) Chars() []rune {
	return self.chars
}
