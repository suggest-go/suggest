package metric

import "math"

// CosineMetric returns a Metric that represents Cosine Metric
func CosineMetric() Metric {
	return &cosine{}
}

type cosine struct{}

func (m *cosine) MinY(alpha float64, size int) int {
	return int(math.Ceil(alpha * alpha * float64(size)))
}

func (m *cosine) MaxY(alpha float64, size int) int {
	return int(math.Floor(float64(size) / (alpha * alpha)))
}

func (m *cosine) Threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(alpha * math.Sqrt(float64(sizeA*sizeB))))
}

func (m *cosine) Distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/math.Sqrt(float64(sizeA*sizeB))
}
