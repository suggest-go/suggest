package suggest

import "math"

type measure interface {
	minY(alpha float64, size int) int
	maxY(alpha float64, size int) int
	threshold(alpha float64, sizeA, sizeB int) int
	distance(inter, sizeA, sizeB int) float64
}

type MeasureT byte

const (
	JACCARD MeasureT = iota
	COSINE
	last
)

var measureHolder [last]measure

func getMeasure(name MeasureT) measure {
	// monkeycode fix me
	if len(measureHolder) <= int(name) || int(name) < 0 {
		panic("Given measure doesn't exists")
	}

	return measureHolder[name]
}

type jaccard struct{}

func (self *jaccard) minY(alpha float64, size int) int {
	return int(math.Ceil(alpha * float64(size)))
}

func (self *jaccard) maxY(alpha float64, size int) int {
	return int(math.Floor(float64(size) / alpha))
}

func (self *jaccard) threshold(alpha float64, sizeA, sizeB int) int {
	return int(math.Ceil(alpha * float64(sizeA+sizeB) / (1 + alpha)))
}

// 1 - |intersection| / |union| = 1 - |intersection| / (|A| + |B| - |intersection|)
func (self *jaccard) distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/float64(sizeA+sizeB-inter)
}

type cosine struct{}

func (self *cosine) minY(alpha float64, size int) int {
	return int(alpha*alpha) * size
}

func (self *cosine) maxY(alpha float64, size int) int {
	return int(float64(size) / (alpha * alpha))
}

func (self *cosine) threshold(alpha float64, sizeA, sizeB int) int {
	return int(alpha * math.Sqrt(float64(sizeA*sizeB)))
}

func (self *cosine) distance(inter, sizeA, sizeB int) float64 {
	return 1 - float64(inter)/math.Sqrt(float64(sizeA*sizeB))
}

func init() {
	measureHolder[JACCARD] = &jaccard{}
	measureHolder[COSINE] = &cosine{}
}
