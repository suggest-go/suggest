package metric

import "math"

// DiceMetric returns a Metric that represents Dice Metric
func DiceMetric() Metric {
	return &dice{}
}

type dice struct{}

func (m *dice) MinY(alpha float64, size int) int {
	return int(math.Ceil(alpha / (2 - alpha) * float64(size)))
}

func (m *dice) MaxY(alpha float64, size int) int {
	return int(math.Floor((2 - alpha) / alpha * float64(size)))
}

func (m *dice) Threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(0.5 * alpha * float64(sizeA+sizeB)))
}

func (m *dice) Distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(2*inter)/float64(sizeA+sizeB)
}
