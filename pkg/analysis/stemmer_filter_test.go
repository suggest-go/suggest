package analysis

import (
	"reflect"
	"testing"
)

func TestRussianStemmerFilter(t *testing.T) {
	filter := NewRussianStemmerFilter()

	actual := filter.Filter(
		[]Token{
			"вместе", "с", "тем", "о", "силе", "электромагнитной", "энергии", "имели", "представление", "еще",
		},
	)

	expected := []Token{
		"вмест", "сил", "электромагнитн", "энерг", "имел", "представлен",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
