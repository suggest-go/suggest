package alphabet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequentialAlphabet(t *testing.T) {
	testCases := []struct {
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

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			alphabet := NewRussianAlphabet()
			assert.Equal(t, testCase.expected, alphabet.Has(testCase.char))
		})
	}
}

func TestCompositeAlphabet(t *testing.T) {
	testCases := []struct {
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

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			alphabet := NewCompositeAlphabet(
				[]Alphabet{
					NewRussianAlphabet(),
					NewEnglishAlphabet(),
					NewNumberAlphabet(),
				},
			)

			assert.Equal(t, testCase.expected, alphabet.Has(testCase.char))
		})
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
