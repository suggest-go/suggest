package metric

// ExactMetric returns a Metric that represents exact matching between 2 ngram sets
func ExactMetric() Metric {
	return &exact{}
}

type exact struct{}

func (m *exact) MinY(alpha float64, size int) int {
	return size
}

func (m *exact) MaxY(alpha float64, size int) int {
	return size
}

func (m *exact) Threshold(alpha float64, sizeA, sizeB int) int {
	return sizeA
}

func (m *exact) Distance(inter, sizeA, sizeB int) float64 {
	return 0
}
