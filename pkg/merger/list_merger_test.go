package merger

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	mergers := []ListMerger{
		ScanCount(),
		CPMerge(),
		MergeSkip(),
		DivideSkip(0.01, MergeSkip()),
	}

	for _, merger := range mergers {
		for _, c := range dataProvider() {
			actual := make(map[int][]uint32, len(c.rid))
			rid := make(Rid, 0, len(c.rid))

			for _, slice := range c.rid {
				rid = append(rid, NewSliceIterator(slice))
			}

			candidates, err := merger.Merge(rid, c.t)

			if err != nil {
				t.Errorf("Unexpected error occurs: %v", err)
			}

			for _, candidate := range candidates {
				actual[candidate.Overlap] = append(actual[candidate.Overlap], candidate.Position)
			}

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
			}
		}
	}
}

type oneCase struct {
	rid      [][]uint32
	t        int
	expected map[int][]uint32
}

func dataProvider() []oneCase {
	return []oneCase{
		{
			[][]uint32{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			2,
			map[int][]uint32{
				2: {1, 3},
				4: {2},
			},
		},
		{
			[][]uint32{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			3,
			map[int][]uint32{
				4: {2},
			},
		},
		{
			[][]uint32{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			4,
			map[int][]uint32{
				4: {2},
			},
		},
		{
			[][]uint32{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			4,
			map[int][]uint32{
				4: {50},
			},
		},
		{
			[][]uint32{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			3,
			map[int][]uint32{
				4: {50},
			},
		},
		{
			[][]uint32{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			2,
			map[int][]uint32{
				2: {10, 60, 100},
				4: {50},
			},
		},
	}
}
