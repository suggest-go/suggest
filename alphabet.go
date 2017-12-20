package suggest

import (
	"sort"
)

// InvalidChar represents unmapped value of given char
const InvalidChar = -1

// Alphabet is abstract for manipulating with set of symbols
type Alphabet interface {
	// MapChar map given char to int32
	MapChar(char rune) int32
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

	return &simpleAlphabet{table, charHolder{chars}}
}

func (a *simpleAlphabet) MapChar(char rune) int32 {
	index, ok := a.table[char]
	if !ok {
		index = InvalidChar
	}

	return index
}

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
		min, max, charHolder{chars},
	}
}

func (a *sequentialAlphabet) MapChar(char rune) int32 {
	if char < a.min || char > a.max {
		return InvalidChar
	}

	return int32(char - a.min)
}

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
func (a *russianAlphabet) MapChar(char rune) int32 {
	if char == 'ё' {
		return a.parent.MapChar('е')
	}

	return a.parent.MapChar(char)
}

func (a *russianAlphabet) Size() int {
	return a.parent.Size()
}

func (a *russianAlphabet) Chars() []rune {
	return a.parent.Chars()
}

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

// compositeAlphabet represents composite pattern for a group of alphabets
type compositeAlphabet struct {
	alphabets []Alphabet
	charHolder
}

// NewCompositeAlphabet returns new instance of compositeAlphabet
func NewCompositeAlphabet(alphabets []Alphabet) Alphabet {
	size := 0

	sort.Slice(alphabets, func(i, j int) bool {
		return alphabets[i].Size() < alphabets[j].Size()
	})

	for _, alphabet := range alphabets {
		size += alphabet.Size()
	}

	chars := make([]rune, 0, size)
	for _, alphabet := range alphabets {
		chars = append(chars, alphabet.Chars()...)
	}

	return &compositeAlphabet{alphabets, charHolder{chars}}
}

// trouble with mapping, we sort alphabet in size order for it, but will be collision when size are equals
func (a *compositeAlphabet) MapChar(char rune) int32 {
	key := int32(InvalidChar)
	shift := 0
	for _, alphabet := range a.alphabets {
		key = alphabet.MapChar(char)
		if key != InvalidChar {
			key += int32(shift)
			break
		}

		shift += alphabet.Size()
	}

	return key
}
