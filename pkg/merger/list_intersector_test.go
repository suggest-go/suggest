package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersect(t *testing.T) {
	testCases := []struct {
		rid      [][]uint32
		expected []uint32
	}{
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

	for _, testCase := range testCases {
		intersector := Intersector()
		rid := make(Rid, 0, len(testCase.rid))

		for _, slice := range testCase.rid {
			rid = append(rid, NewSliceIterator(slice))
		}

		collector := &SimpleCollector{}
		assert.NoError(t, intersector.Intersect(rid, collector))

		actual := []uint32{}

		for _, c := range collector.Candidates {
			actual = append(actual, c.Position())
		}

		assert.Equal(t, testCase.expected, actual)
	}
}
