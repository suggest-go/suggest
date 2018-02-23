package metric

import "math"

// OverlapMetric returns a Metric that represents Overlap metric
func OverlapMetric() Metric {
	return &overlap{}
}

type overlap struct{}

func (m *overlap) MinY(alpha float64, size int) int {
	return 1
}

func (m *overlap) MaxY(alpha float64, size int) int {
	return math.MaxInt16
}

func (m *overlap) Threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(alpha * math.Min(float64(sizeA), float64(sizeB))))
}

func (m *overlap) Distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/(math.Min(float64(sizeA), float64(sizeB)))
}
