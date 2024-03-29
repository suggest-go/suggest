package merger

// ScanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
func ScanCount() ListMerger {
	return newMerger(&scanCount{})
}

type scanCount struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (lm *scanCount) Merge(rid Rid, threshold int, collector Collector) error {
	size := len(rid)
	candidates := make([]MergeCandidate, 0, size)
	tmp := make([]MergeCandidate, 0, size)
	j, endMergeCandidate := 0, 0

	for _, list := range rid {
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

	tmp = tmp[:0]

	for _, c := range candidates {
		if c.Overlap() >= threshold {
			err := collector.Collect(c)

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
