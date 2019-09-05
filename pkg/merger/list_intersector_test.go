package merger

import (
	"reflect"
	"testing"
)

func TestIntersect(t *testing.T) {
	intersectors := []ListIntersector{
		Intersector(),
	}

	for _, intersector := range intersectors {
		for _, c := range intersectDataProvider() {
			rid := make(Rid, 0, len(c.rid))

			for _, slice := range c.rid {
				rid = append(rid, NewSliceIterator(slice))
			}

			collector := &SimpleCollector{}
			err := intersector.Intersect(rid, collector)

			if err != nil {
				t.Errorf("Unexpected error occurs: %v", err)
			}

			actual := []uint32{}

			for _, c := range collector.Candidates {
				actual = append(actual, c.Position())
			}

			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("Test Fail, expected %v, got %v", c.expected, actual)
			}
		}
	}
}

type intersectCase struct {
	rid      [][]uint32
	expected []uint32
}

func intersectDataProvider() []intersectCase {
	return []intersectCase{
		{
			[][]uint32{
				{1, 2, 3},
				{1, 2},
				{2, 3},
				{2},
			},
			[]uint32{
				2,
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
			[]uint32{},
		},
		{
			[][]uint32{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{50, 100, 200},
			},
			[]uint32{
				50,
			},
		},
	}
}
