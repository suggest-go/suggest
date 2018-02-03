package suggest

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	mergers := []ListMerger{
		&ScanCount{},
		&CPMerge{},
		&MergeSkip{},
		&DivideSkip{0.01, &MergeSkip{}},
	}

	for _, merger := range mergers {
		for _, c := range dataProvider() {
			actual := make(map[int]PostingList, len(c.rid))

			for _, candidate := range merger.Merge(c.rid, c.t) {
				if candidate.Overlap >= c.t {
					actual[candidate.Overlap] = append(actual[candidate.Overlap], candidate.Pos)
				}
			}

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
			}
		}
	}
}

type oneCase struct {
	rid      Rid
	t        int
	expected map[int]PostingList
}

func dataProvider() []oneCase {
	return []oneCase{
		{
			Rid{
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
			Rid{
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
			Rid{
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
			Rid{
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
			Rid{
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
			Rid{
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
