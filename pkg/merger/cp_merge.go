package merger

import "sort"

// CPMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
// inspired by https://github.com/chokkan/simstring
func CPMerge() ListMerger {
	return &cpMerge{}
}

type cpMerge struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (cp *cpMerge) Merge(rid Rid, threshold int) []*MergeCandidate {
	lenRid := len(rid)
	minQueries := lenRid - threshold + 1
	candidates := make([]*MergeCandidate, 0, lenRid)

	if threshold > lenRid {
		return candidates
	}

	sort.Sort(rid)

	tmp := make([]*MergeCandidate, 0, lenRid)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid[:minQueries] {
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

	if len(candidates) == 0 {
		return candidates
	}

	for i := minQueries; i < lenRid; i++ {
		tmp = tmp[:0]

		for _, c := range candidates {
			j := lowerBound(rid[i], c.Position)

			if j != -1 {
				if rid[i][j] == c.Position {
					c.Overlap++
				}

				rid[i] = rid[i][j:]
			}

			if c.Overlap+(lenRid-i-1) >= threshold {
				tmp = append(tmp, c)
			}
		}

		candidates, tmp = tmp, candidates

		if len(candidates) == 0 {
			break
		}
	}

	return candidates
}