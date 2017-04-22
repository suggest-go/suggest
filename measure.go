package suggest

import "math"

type measure interface {
	minY(alpha float64, cardinality int) int
	maxY(alpha float64, cardinality int) int
	threshold(alpha float64, aCardinality, bCardinality int) int
}

type jaccard struct{}

func (self *jaccard) minY(alpha float64, cardinality int) int {
	return int(math.Ceil(alpha * float64(cardinality)))
}

func (self *jaccard) maxY(alpha float64, cardinality int) int {
	return int(math.Floor(float64(cardinality) / alpha))
}

func (self *jaccard) threshold(alpha float64, aCardinality, bCardinality int) int {
	return int(math.Ceil(alpha * float64(aCardinality+bCardinality) / (1 + alpha)))
}
