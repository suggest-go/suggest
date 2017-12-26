package suggest

import "testing"

func TestSearchUpperBound(t *testing.T) {
	items := PostingList{0, 1, 3, 7, 9, 10, 11}
	cases := []struct {
		val      Position
		expected int
	}{
		{1, 1},
		{13, -1},
		{2, 3},
		{5, 7},
		{10, 10},
		{9, 9},
		{8, 9},
		{7, 7},
		{6, 7},
	}

	for _, c := range cases {
		actual := binarySearchUpperBound(items, c.val)
		if actual != -1 {
			actual = int(items[actual])
		}

		if actual != c.expected {
			t.Errorf("Test Fail, expected %d, got %d", c.expected, actual)
		}
	}
}

func TestBinarySearch(t *testing.T) {
	items := PostingList{0, 1, 3, 7, 9, 10, 11}
	cases := []struct {
		val      Position
		expected int
	}{
		{1, 1},
		{13, -1},
		{2, -1},
		{5, -1},
		{10, 10},
		{9, 9},
		{8, -1},
		{7, 7},
		{6, -1},
	}

	for _, c := range cases {
		actual := binarySearch(items, c.val)
		if actual != -1 {
			actual = int(items[actual])
		}

		if actual != c.expected {
			t.Errorf("Test Fail, expected %d, got %d", c.expected, actual)
		}
	}
}
