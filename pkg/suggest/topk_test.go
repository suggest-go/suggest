package suggest

import (
	"reflect"
	"testing"
)

func TestTopKQueue(t *testing.T) {
	candidates := []Candidate{
		{Key: 1, Score: 0.1},
		{Key: 2, Score: 0.01},
		{Key: 3, Score: 0.91},
		{Key: 4, Score: 0.24},
		{Key: 5, Score: 0.13},
		{Key: 6, Score: 0.07},
		{Key: 7, Score: 0.9},
		{Key: 8, Score: 0.12345},
		{Key: 9, Score: 0.65},
		{Key: 10, Score: 0.6565},
	}

	selector := NewTopKQueue(3)

	for _, candidate := range candidates {
		selector.Add(candidate.Key, candidate.Score)
	}

	expected := []Candidate{
		{Key: 3, Score: 0.91},
		{Key: 7, Score: 0.9},
		{Key: 10, Score: 0.6565},
	}

	actual := selector.GetCandidates()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Test fail, expected %v, got %v", expected, actual)
	}

	if selector.GetLowestScore() != 0.6565 {
		t.Errorf("Test fail, the lowest score should be %v, got %v", 0.6565, selector.GetLowestScore())
	}

	if !selector.CanTakeWithScore(0.6566) {
		t.Errorf("Test fail, CanTakeWithScore should return true")
	}
}
