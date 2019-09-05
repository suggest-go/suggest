package merger

import "sort"

// ListIntersector is the interface that is responsible for intersection operation
// between array of docs iterators
type ListIntersector interface {
	// Intersect performs intersection operation for the given rid and
	// transmits the result to collector
	Intersect(rid Rid, collector Collector) error
}

// intersector implements ListIntersector interface
type intersector struct{}

// Intersector creates a new instance of ListIntersector
func Intersector() ListIntersector {
	return &intersector{}
}

// Intersect performs intersection operation for the given rid and
// transmits the result to collector
func (i *intersector) Intersect(rid Rid, collector Collector) error {
	n := len(rid)

	if n == 0 {
		return nil
	}

	sort.Sort(rid)
	first := rid[0]
	rest := rid[1:]

	item, err := first.Get()

	if err != nil {
		return err
	}

	for {
		isGoodCandidate := true

		for _, it := range rest {
			lower, err := it.LowerBound(item)

			if err == ErrIteratorIsNotDereferencable || (err == nil && lower != item) {
				isGoodCandidate = false
				break
			}

			if err != nil {
				return err
			}
		}

		if isGoodCandidate {
			err := collector.Collect(NewMergeCandidate(item, n))

			// we are going to ignore it
			if err == ErrCollectionTerminated {
				return nil
			}

			if err != nil {
				return err
			}
		}

		if !first.HasNext() {
			break
		}

		item, err = first.Next()

		if err != nil {
			return err
		}
	}

	return nil
}
