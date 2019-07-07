package merger

import (
	"math"
	"sort"
)

// DivideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
func DivideSkip(mu float64, merger ListMerger) ListMerger {
	return &divideSkip{
		mu:     mu,
		merger: merger,
	}
}

type divideSkip struct {
	mu     float64
	merger ListMerger
}

// Merge returns list of candidates, that appears at least `threshold` times.
func (ds *divideSkip) Merge(rid Rid, threshold int) ([]MergeCandidate, error) {
	sort.Sort(sort.Reverse(rid))

	M := float64(rid[0].Len())
	l := int(float64(threshold) / (ds.mu*math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]

	if len(lShort) == 0 {
		return ds.merger.Merge(rid, threshold)
	}

	candidates, err := ds.merger.Merge(lShort, threshold-l)

	if err != nil {
		return nil, err
	}

	var (
		position   uint32
		result     = make([]MergeCandidate, 0, len(candidates))
	)

	for _, c := range candidates {
		position = c.Position

		for _, longList := range lLong {
			r, err := longList.LowerBound(position)

			if err != nil && err != ErrIteratorIsNotDereferencable {
				return nil, err
			}

			if err == nil && r == position {
				c.Overlap++
			}
		}

		if c.Overlap >= threshold {
			result = append(result, c)
		}
	}

	return result, nil
}
