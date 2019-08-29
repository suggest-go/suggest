package merger

import (
	"sort"
	"sync"
)

// CPMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
// inspired by https://github.com/chokkan/simstring
func CPMerge() ListMerger {
	return newMerger(&cpMerge{})
}

type cpMerge struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (cp *cpMerge) Merge(rid Rid, threshold int, collector Collector) error {
	lenRid := len(rid)
	minQueries := lenRid - threshold + 1
	j, endMergeCandidate := 0, 0

	sort.Sort(rid)

	tmp := bufPool.Get().([]MergeCandidate)
	candidates := bufPool.Get().([]MergeCandidate)

	for _, list := range rid[:minQueries] {
		isValid := true
		current, err := list.Get()

		if err != nil {
			if err == ErrIteratorIsNotDereferencable {
				isValid = false
			} else {
				return err
			}
		}

		tmp = tmp[:0]
		j, endMergeCandidate = 0, len(candidates)

		for j < endMergeCandidate || isValid {
			if j >= endMergeCandidate || (isValid && candidates[j].Position() > current) {
				tmp = append(tmp, NewMergeCandidate(current, 1))

				if list.HasNext() {
					current, err = list.Next()

					if err != nil {
						return err
					}
				} else {
					isValid = false
				}
			} else if !isValid || (j < endMergeCandidate && candidates[j].Position() < current) {
				tmp = append(tmp, candidates[j])
				j++
			} else {
				candidates[j].increment()
				tmp = append(tmp, candidates[j])
				j++

				if list.HasNext() {
					current, err = list.Next()

					if err != nil {
						return err
					}
				} else {
					isValid = false
				}
			}
		}

		candidates, tmp = tmp, candidates
	}

	for i := minQueries; i < lenRid; i++ {
		tmp = tmp[:0]

		for _, c := range candidates {
			current, err := rid[i].LowerBound(c.Position())

			if err == nil && current == c.Position() {
				c.increment()
			}

			if err != nil && err != ErrIteratorIsNotDereferencable {
				return err
			}

			if lenRid-1 == i && c.Overlap() >= threshold {
				err = collector.Collect(c)

				if err == ErrCollectionTerminated {
					return nil
				}

				if err != nil {
					return err
				}

				continue
			}

			if c.Overlap()+(lenRid-i-1) >= threshold {
				tmp = append(tmp, c)
			}
		}

		candidates, tmp = tmp, candidates

		if len(candidates) == 0 {
			break
		}
	}

	bufPool.Put(tmp[:0])
	bufPool.Put(candidates[:0])

	return nil
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]MergeCandidate, 0, 1024)
	},
}
