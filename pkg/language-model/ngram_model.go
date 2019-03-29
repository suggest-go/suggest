package lm

import (
	"errors"
	"math"
)

//
type NGramModel interface {
	//
	Put(nGrams []WordID, count WordCount) error
	//
	Score(nGrams []WordID) float64
	//
	Next(nGrams []WordID) ([]WordID, error)
}

const (
	alpha            = 0.4
	unknownWordScore = -100.0
)

type nGramModel struct {
	indices    []NGramVector
	nGramOrder uint8
}

//
func NewNGramModel(nGramOrder uint8) NGramModel {
	indices := make([]NGramVector, nGramOrder)

	for i := 0; i < int(nGramOrder); i++ {
		indices[i] = NewNGramVector()
	}

	return &nGramModel{
		indices:    indices,
		nGramOrder: nGramOrder,
	}
}

//
func (m *nGramModel) Put(nGrams []WordID, count WordCount) error {
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

//
func (m *nGramModel) Score(nGrams []WordID) float64 {
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

//
func (m *nGramModel) Next(nGrams []WordID) ([]WordID, error) {
	if int(m.nGramOrder) <= len(nGrams) {
		return nil, errors.New("nGrams length should be less than the nGramModel order")
	}

	order := 0
	count := WordCount(0)
	parent := InvalidContextOffset

	for ; order < len(nGrams); order++ {
		vector := m.indices[order]
		count, parent = vector.GetCount(nGrams[order], parent)

		if count == 0 {
			return []WordID{}, nil
		}
	}

	return m.indices[order].Next(parent), nil
}

//
func calcScore(counts []WordCount) float64 {
	factor := float64(1)

	for i := len(counts) - 1; i >= 1; i-- {
		if counts[i] > 0 {
			return math.Log(factor * float64(counts[i]) / float64(counts[i-1]))
		}

		factor *= alpha
	}

	return unknownWordScore
}
