package merger

import (
	"math"
	"sort"
)

// DivideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
func DivideSkip(mu float64) ListMerger {
	return newMerger(&divideSkip{
		mu:     mu,
		merger: MergeSkip(),
	})
}

type divideSkip struct {
	mu     float64
	merger ListMerger
}

// Merge returns list of candidates, that appears at least `threshold` times.
func (ds *divideSkip) Merge(rid Rid, threshold int, collector Collector) error {
	sort.Sort(sort.Reverse(rid))

	M := float64(rid[0].Len())
	l := int(float64(threshold) / (ds.mu*math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]

	if len(lShort) == 0 {
		return ds.merger.Merge(rid, threshold, collector)
	}

	mergeRes := &SimpleCollector{}
	err := ds.merger.Merge(lShort, threshold-l, mergeRes)

	if err != nil {
		return err
	}

	for _, c := range mergeRes.Candidates {
		position := c.Position()

		for _, longList := range lLong {
			r, err := longList.LowerBound(position)

			if err != nil && err != ErrIteratorIsNotDereferencable {
				return err
			}

			if err == nil && r == position {
				c.increment()
			}
		}

		if c.Overlap() >= threshold {
			err = collector.Collect(c)

			if err == ErrCollectionTerminated {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}
