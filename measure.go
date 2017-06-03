package suggest

import "math"

// Metric defined here, is not pure mathematics metric definition as distance between each pair of elements of a set.
// Here we can also ask metric to give as minimum intersection between A and B for given alpha,
// min/max candidate cardinality
type Metric interface {
	// MinY returns the minimum ngram cardinality for candidate
	MinY(alpha float64, size int) int
	// MinY returns the maximum ngram cardinality for candidate
	MaxY(alpha float64, size int) int
	// Threshold returns required intersection between A and B for given alpha
	Threshold(alpha float64, sizeA, sizeB int) int
	// Distance calculate distance between 2 sets
	Distance(inter, sizeA, sizeB int) float64
}

type jaccard struct{}

// JaccardMetric returns a Metric that represents Jaccard Metric
func JaccardMetric() Metric {
	return &jaccard{}
}

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
