package analysis

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRussianStemmerFilter(t *testing.T) {
	testCases := []struct {
		sentence string
		expected []Token
	}{
		{
			sentence: "вместе с тем о силе электромагнитной энергии имели представление еще",
			expected: []Token{"вмест", "сил", "электромагнитн", "энерг", "имел", "представлен"},
		},
		{
			sentence: "total 2310 рублей итого",
			expected: []Token{"total", "2310", "рубл", "ит"},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			filter := NewRussianStemmerFilter()

			actual := filter.Filter(
				strings.Split(testCase.sentence, " "),
			)

			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected %v, got %v", testCase.expected, actual)
			}
		})
	}
}

func TestEnglishStemmerFilter(t *testing.T) {
	testCases := []struct {
		sentence string
		expected []Token
	}{
		{
			sentence: "What does борщ mean",
			expected: []Token{"What", "борщ", "mean"},
		},
		{
			sentence: "Hello hello mister Credo What's up",
			expected: []Token{"Hello", "hello", "mister", "Credo", "What"},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			filter := NewEnglishStemmerFilter()

			actual := filter.Filter(
				strings.Split(testCase.sentence, " "),
			)

			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected %v, got %v", testCase.expected, actual)
			}
		})
	}
}
