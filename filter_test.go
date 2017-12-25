package suggest

import (
	"reflect"
	"testing"
	"sort"
)

// IMPLEMENT ME
func TestCPMerge(t *testing.T) {
	for _, c := range dataProvider() {
		actual := cpMerge(c.rid, c.t)
		actualMap := make(map[int]PostingList)
		for n, list := range actual {
			if len(list) == 0 {
				continue
			}

			actualMap[n] = list
		}

		if !reflect.DeepEqual(actualMap, c.expected) {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}
}

func TestScanCount(t *testing.T) {
	for _, c := range dataProvider() {
		actual := scanCount(c.rid, c.t)
		actualMap := make(map[int]PostingList)
		for n, list := range actual {
			if len(list) == 0 {
				continue
			}

			sort.Slice(list, func(i, j int) bool {
				return list[i] < list[j]
			})

			actualMap[n] = list
		}

		if !reflect.DeepEqual(actualMap, c.expected) {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}
}

func TestDivideSkip(t *testing.T) {
	for _, c := range dataProvider() {
		actual := divideSkip(c.rid, c.t, 0.0085)
		actualMap := make(map[int]PostingList)
		for n, list := range actual {
			if len(list) == 0 {
				continue
			}

			actualMap[n] = list
		}

		if !reflect.DeepEqual(actualMap, c.expected) {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}
}

func TestMergeSkip(t *testing.T) {
	for _, c := range dataProvider() {
		actual := mergeSkip(c.rid, c.t)
		actualMap := make(map[int]PostingList)
		for n, list := range actual {
			if len(list) == 0 {
				continue
			}

			actualMap[n] = list
		}

		if !reflect.DeepEqual(actualMap, c.expected) {
			t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
		}
	}
}

type oneCase struct {
	rid      []PostingList
	t        int
	expected map[int]PostingList
}

func dataProvider() []oneCase {
	return []oneCase{
		{
			[]PostingList{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			2,
			map[int]PostingList{
				2: {1, 3},
				4: {2},
			},
		},
		{
			[]PostingList{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			3,
			map[int]PostingList{
				4: {2},
			},
		},
		{
			[]PostingList{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			4,
			map[int]PostingList{
				4: {2},
			},
		},
		{
			[]PostingList{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			4,
			map[int]PostingList{
				4: {50},
			},
		},
		{
			[]PostingList{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			3,
			map[int]PostingList{
				4: {50},
			},
		},
		{
			[]PostingList{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			2,
			map[int]PostingList{
				2: {10, 60, 100},
				4: {50},
			},
		},
	}
}
