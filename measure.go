package suggest

import "math"

// measure
type measure interface {
	// minY returns the minimum ngram cardinality for candidate
	minY(alpha float64, size int) int
	// minY returns the maximum ngram cardinality for candidate
	maxY(alpha float64, size int) int
	// threshold returns required intersection between A and B for given alpha
	threshold(alpha float64, sizeA, sizeB int) int
	// distance calculate distance between 2 strings
	distance(inter, sizeA, sizeB int) float64
}

// MeasureT represents type of measure name
type MeasureT byte

const (
	// Jaccard represents Jaccard distance
	Jaccard MeasureT = iota
	// Cosine represents Cosine distance
	Cosine
	// Dice represents Dice distance
	Dice
	// Exact represents ngrams sets equality
	Exact
	last
)

// measureHolder is a holder for existing measures
var measureHolder [last]measure

// getMeasure returns measure implements by given name
func getMeasure(name MeasureT) measure {
	// monkeycode fix me
	if len(measureHolder) <= int(name) || int(name) < 0 {
		panic("Given measure doesn't exists")
	}

	return measureHolder[name]
}

type jaccard struct{}

func (m *jaccard) minY(alpha float64, size int) int {
	return int(math.Ceil(alpha * float64(size)))
}

func (m *jaccard) maxY(alpha float64, size int) int {
	return int(math.Floor(float64(size) / alpha))
}

func (m *jaccard) threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(alpha * float64(sizeA+sizeB) / (1 + alpha)))
}

// 1 - |intersection| / |union| = 1 - |intersection| / (|A| + |B| - |intersection|)
func (m *jaccard) distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/float64(sizeA+sizeB-inter)
}

type cosine struct{}

func (m *cosine) minY(alpha float64, size int) int {
	return int(alpha * alpha * float64(size))
}

func (m *cosine) maxY(alpha float64, size int) int {
	return int(float64(size) / (alpha * alpha))
}

func (m *cosine) threshold(alpha float64, sizeA, sizeB int) int {
	return int(alpha * math.Sqrt(float64(sizeA*sizeB)))
}

func (m *cosine) distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/math.Sqrt(float64(sizeA*sizeB))
}

type dice struct{}

func (m *dice) minY(alpha float64, size int) int {
	return int(alpha / (2 - alpha) * float64(size))
}

func (m *dice) maxY(alpha float64, size int) int {
	return int((2 - alpha) / alpha * float64(size))
}

func (m *dice) threshold(alpha float64, sizeA, sizeB int) int {
	return int(0.5 * alpha * float64(sizeA+sizeB))
}

func (m *dice) distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(2*inter)/float64(sizeA+sizeB)
}

type exact struct{}

func (m *exact) minY(alpha float64, size int) int {
	return size
}

func (m *exact) maxY(alpha float64, size int) int {
	return size
}

func (m *exact) threshold(alpha float64, sizeA, sizeB int) int {
	return sizeA
}

func (m *exact) distance(inter, sizeA, sizeB int) float64 {
	return 0
}

func init() {
	measureHolder[Jaccard] = &jaccard{}
	measureHolder[Cosine] = &cosine{}
	measureHolder[Dice] = &dice{}
	measureHolder[Exact] = &exact{}
}
