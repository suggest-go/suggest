package alphabet

import (
	"testing"
)

func TestSequentialAlphabet(t *testing.T) {
	alphabet := NewRussianAlphabet()

	cases := []struct {
		char     rune
		expected bool
	}{
		{'а', true},
		{'е', true},
		{'ё', true},
		{'я', true},
		{'j', false},
		{'7', false},
	}

	for _, c := range cases {
		actual := alphabet.Has(c.char)

		if c.expected != actual {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
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
		expected bool
	}{
		{'a', true},
		{'b', true},
		{'z', true},
		{'а', true},
		{'ё', true},
		{'е', true},
		{'ж', true},
		{'я', true},
		{'7', true},
		{'-', false},
	}

	for _, c := range cases {
		actual := alphabet.Has(c.char)
		if c.expected != actual {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}
}

func BenchmarkHas(b *testing.B) {
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
			alphabet.Has(runeVal)
		}
	}
}
