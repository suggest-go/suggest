package metric

import "math"

// JaccardMetric returns a Metric that represents Jaccard Metric
func JaccardMetric() Metric {
	return &jaccard{}
}

type jaccard struct{}

func (m *jaccard) MinY(alpha float64, size int) int {
	return int(math.Ceil(alpha * float64(size)))
}

func (m *jaccard) MaxY(alpha float64, size int) int {
	return int(math.Floor(float64(size) / alpha))
}

func (m *jaccard) Threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(alpha * float64(sizeA+sizeB) / (1 + alpha)))
}

// 1 - |intersection| / |union| = 1 - |intersection| / (|A| + |B| - |intersection|)
func (m *jaccard) Distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/float64(sizeA+sizeB-inter)
}
