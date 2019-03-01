package language_model

import (
	"errors"
)

const alpha = 0.4

type nGramModel struct {
	indices    []*sortedArray
	nGramOrder uint8
}

func NewNGramModel(nGramOrder uint8) *nGramModel {
	indices := make([]*sortedArray, nGramOrder)

	for i := 0; i < int(nGramOrder); i++ {
		indices[i] = NewNGramVector()
	}

	return &nGramModel{
		indices:    indices,
		nGramOrder: nGramOrder,
	}
}

func (m *nGramModel) Put(nGrams []WordId, count WordId) error {
	if len(nGrams) > int(m.nGramOrder) {
		return errors.New("nGrams order is out of range")
	}

	parent := InvalidContextOffset

	for i := 0; i < len(nGrams); i++ {
		if i == len(nGrams)-1 {
			m.indices[i].Put(nGrams[i], parent, count)
		} else {
			parent = m.indices[i].GetContextOffset(nGrams[i], parent)
		}
	}

	return nil
}

func (m *nGramModel) Score(nGrams []WordId) float64 {
	order := int(m.nGramOrder)
	if order > len(nGrams) {
		order = len(nGrams)
	}

	counts := make([]WordCount, order+1)
	parent := InvalidContextOffset

	for i := 0; i < order; i++ {
		vector := m.indices[i]

		if i == 0 {
			counts[0] = vector.CorpousCount()
		}

		counts[i+1], parent = vector.GetCount(nGrams[i], parent)
	}

	return calcScore(counts)
}

func calcScore(counts []WordCount) float64 {
	factor := float64(1)

	for i := len(counts) - 1; i >= 1; i-- {
		if counts[i] > 0 {
			return factor * float64(counts[i]) / float64(counts[i-1])
		}

		factor *= alpha
	}

	return 0.0
}
