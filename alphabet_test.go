package suggest

import (
	"testing"
)

func TestSequentialAlphabet(t *testing.T) {
	alphabet := NewRussianAlphabet()

	cases := []struct {
		char     rune
		expected int32
	}{
		{'а', 0},
		{'е', 5},
		{'ё', 5},
		{'я', 31},
		{'j', InvalidChar},
		{'7', InvalidChar},
	}

	for _, c := range cases {
		actual := alphabet.MapChar(c.char)
		if c.expected != actual {
			t.Errorf("Test Fail, expected %d, got %d", c.expected, actual)
		}
	}
}

func TestCompositeAlphabet(t *testing.T) {
	alphabet := NewCompositeAlphabet(
		[]Alphabet{
			NewRussianAlphabet(),
			NewEnglishAlphabet(),
			NewNumberAlphabet(),
		},
	)

	cases := []struct {
		char     rune
		expected int32
	}{
		{'a', 10 + 0},
		{'b', 10 + 1},
		{'z', 10 + 25},
		{'а', 10 + 26 + 0},
		{'ё', 10 + 26 + 5},
		{'е', 10 + 26 + 5},
		{'ж', 10 + 26 + 6},
		{'я', 10 + 26 + 31},
		{'7', 7},
		{'-', InvalidChar},
	}

	for _, c := range cases {
		actual := alphabet.MapChar(c.char)
		if c.expected != actual {
			t.Errorf("Test Fail, expected %d, got %d", c.expected, actual)
		}
	}
}

func BenchmarkMapChar(b *testing.B) {
	ngram := "ёj9"
	alphabet := NewCompositeAlphabet(
		[]Alphabet{
			NewEnglishAlphabet(),
			NewRussianAlphabet(),
			NewNumberAlphabet(),
		},
	)

	for i := 0; i < b.N; i++ {
		for _, runeVal := range ngram {
			alphabet.MapChar(runeVal)
		}
	}
}
