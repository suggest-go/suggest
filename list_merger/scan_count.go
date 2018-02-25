package list_merger

// ScanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
func ScanCount() ListMerger {
	return &scanCount{}
}

type scanCount struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (lm *scanCount) Merge(rid Rid, threshold int) []*MergeCandidate {
	size := len(rid)
	candidates := make([]*MergeCandidate, 0, size)
	tmp := make([]*MergeCandidate, 0, size)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid {
		j, k = 0, 0
		tmp = tmp[:0]
		endMergeCandidate, endRid = len(candidates), len(list)

		for j < endMergeCandidate || k < endRid {
			if j >= endMergeCandidate || (k < endRid && candidates[j].Position > list[k]) {
				tmp = append(tmp, &MergeCandidate{list[k], 1})
				k++
			} else if k >= endRid || (j < endMergeCandidate && candidates[j].Position < list[k]) {
				tmp = append(tmp, candidates[j])
				j++
			} else {
				candidates[j].Overlap++
				tmp = append(tmp, candidates[j])
				j++
				k++
			}
		}

		candidates, tmp = tmp, candidates
	}

	tmp = tmp[:0]

	for _, c := range candidates {
		if c.Overlap >= threshold {
			tmp = append(tmp, c)
		}
	}

	candidates = tmp
	return candidates
}
