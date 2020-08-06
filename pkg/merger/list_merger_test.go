package merger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/utils"
)

func TestMergeOverlapOverflow(t *testing.T) {
	m := NewMergeCandidate(1, MaxOverlap)

	assert.Panics(t, func() {
		m.increment()
	})
}

func BenchmarkMergeCandidate(b *testing.B) {
	m := NewMergeCandidate(1, 1)
	p, o := uint32(0), int(0)
	c := utils.Max(100, b.N/MaxOverlap+1)
	e := int(1)

	for i := 0; i < b.N; i++ {
		if i%c == 0 {
			m.increment()
			e++
		}

		for j := 0; j < 100; j++ {
			p = m.Position()
			o = m.Overlap()
		}
	}

	if p != 1 || o != e {
		b.Errorf("Test fail, expected p = 1 && o = %d, got p = %d && o = %d", e, p, o)
	}
}

func TestMerge(t *testing.T) {
	var testCases = []struct {
		rid      [][]uint32
		t        int
		expected map[int][]uint32
	}{
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
		// issue#28
		{
			[][]uint32{
				{1, 2, 3, 5, 7, 10, 30, 50},
				{10, 11, 13, 16, 50, 60, 131},
				{40, 50, 60},
				{50, 100},
				{100, 200},
			},
			1,
			map[int][]uint32{
				1: {1, 2, 3, 5, 7, 11, 13, 16, 30, 40, 131, 200},
				2: {10, 60, 100},
				4: {50},
			},
		},
	}

	mergers := []struct {
		name   string
		merger ListMerger
	}{
		{"scan_count", ScanCount()},
		{"cp_merge", CPMerge()},
		{"merge_skip", MergeSkip()},
		{"divide_skip", DivideSkip(0.01)},
	}

	for _, data := range mergers {
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("Test #%d, merger: %s", i+1, data.name), func(t *testing.T) {
				actual := make(map[int][]uint32, len(testCase.rid))
				rid := make(Rid, 0, len(testCase.rid))

				for _, slice := range testCase.rid {
					rid = append(rid, NewSliceIterator(slice))
				}

				collector := &SimpleCollector{}
				err := data.merger.Merge(rid, testCase.t, collector)
				assert.NoError(t, err)

				for _, candidate := range collector.Candidates {
					actual[candidate.Overlap()] = append(actual[candidate.Overlap()], candidate.Position())
				}

				assert.Equal(t, testCase.expected, actual)
			})
		}
	}
}
