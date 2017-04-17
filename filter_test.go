package suggest

import (
	"reflect"
	"testing"
)

func TestDivideSkip(t *testing.T) {
	divideSkip(
		[][]int{
			{1, 2, 3, 5, 7, 10, 30, 50},
			{10, 11, 13, 16, 50, 60, 131},
			{40, 50, 60},
			{50, 100},
			{100, 200},
		},
		2,
	)
}

func TestMergeSkip(t *testing.T) {
	cases := []struct {
		rid      [][]int
		t        int
		expected map[int][]int
	}{
		{
			[][]int{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			2,
			map[int][]int{
				2: []int{1, 3},
				4: []int{2},
			},
		},
		{
			[][]int{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			3,
			map[int][]int{
				4: []int{2},
			},
		},
		{
			[][]int{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			4,
			map[int][]int{
				4: []int{2},
			},
		},
		{
			[][]int{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			4,
			map[int][]int{
				4: []int{50},
			},
		},
		{
			[][]int{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			3,
			map[int][]int{
				4: []int{50},
			},
		},
		{
			[][]int{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			2,
			map[int][]int{
				2: []int{10, 60, 100},
				4: []int{50},
			},
		},
	}

	for _, c := range cases {
		actual := mergeSkip(c.rid, c.t)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}

}

func TestBinSearch(t *testing.T) {
	items := []int{0, 1, 3, 7, 9, 10, 11}
	cases := []struct {
		val      int
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
		actual := binarySearch(items, 0, c.val)
		if actual != -1 {
			actual = items[actual]
		}

		if actual != c.expected {
			t.Errorf("Test Fail, expected %d, got %d", c.expected, actual)
		}
	}
}
