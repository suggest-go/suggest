package lm

import (
	"errors"
	"fmt"
	"math"

	"github.com/suggest-go/suggest/pkg/store"
)

// NGramModel is an entity that responses for scoring the given nGrams
type NGramModel interface {
	store.Marshaler
	store.Unmarshaler

	// Score returns a lm value of the given sequence of WordID
	Score(nGrams []WordID) float64
	// Next returns a list of WordID which follow after the given sequence of nGrams
	Next(nGrams []WordID) (ScorerNext, error)
}

const (
	// UnknownWordScore is the score for unknown phrases
	UnknownWordScore = -100.0
	alpha            = 0.4
	modelVersion     = "0.0.2"
)

// nGramModel implements NGramModel Stupid backoff
type nGramModel struct {
	indices    []NGramVector
	nGramOrder uint8
}

// NewNGramModel creates a new empty instance of NGramModel instance.
func NewNGramModel() NGramModel {
	return &nGramModel{}
}

// CreateNGramModel creates a NGramModel from the given indices.
func CreateNGramModel(indices []NGramVector) NGramModel {
	return &nGramModel{
		indices:    indices,
		nGramOrder: uint8(len(indices)),
	}
}

// Score returns a lm value of the given sequence of WordID
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
			counts[0] = vector.CorpusCount()
		}

		counts[i+1], parent = vector.GetCount(nGrams[i], parent)
	}

	return calcScore(counts)
}

// Next returns a list of WordID where each candidate follows after the given sequence of nGrams
func (m *nGramModel) Next(nGrams []WordID) (ScorerNext, error) {
	if int(m.nGramOrder) <= len(nGrams) || len(nGrams) == 0 {
		return nil, errors.New("nGrams length should be less than the nGramModel order")
	}

	order := 0
	counts := make([]WordCount, 0, len(nGrams))
	count := WordCount(0)
	parent := InvalidContextOffset

	for ; order < len(nGrams); order++ {
		vector := m.indices[order]
		count, parent = vector.GetCount(nGrams[order], parent)

		if count == 0 {
			return nil, nil
		}

		counts = append(counts, count)
	}

	subVector := m.indices[order].SubVector(parent)

	if subVector == nil {
		return nil, nil
	}

	return &scorerNext{
		contextCounts: counts,
		nGramVector:   subVector,
		context:       parent,
	}, nil
}

func (m *nGramModel) Store(out store.Output) (int, error) {
	if n, err := out.Write([]byte(modelVersion)); err != nil {
		return n, err
	}

	if err := out.WriteByte(byte(m.nGramOrder)); err != nil {
		return 0, err
	}

	p := 6 // 5 + 1

	for _, vector := range m.indices {
		v := vector.(*packedArray) // TODO fix me
		n, err := v.Store(out)
		p += n

		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func (m *nGramModel) Load(in store.Input) (int, error) {
	version := make([]byte, 5)
	n, err := in.Read(version)
	p := n

	if err != nil {
		return p, err
	}

	if string(version) != modelVersion {
		return p, fmt.Errorf("Version mismatch, expected %s, got %s", modelVersion, version)
	}

	order, err := in.ReadByte()
	p++

	if err != nil {
		return p, err
	}

	m.nGramOrder = uint8(order)
	m.indices = make([]NGramVector, m.nGramOrder)

	for i := uint8(0); i < m.nGramOrder; i++ {
		vector := NewNGramVector()
		n, err := vector.Load(in)
		p += n

		if err != nil {
			return p, err
		}

		m.indices[i] = vector
	}

	return p, nil
}

// calcScore returns score for the given counts
func calcScore(counts []WordCount) float64 {
	factor := float64(1)

	for i := len(counts) - 1; i >= 1; i-- {
		if counts[i] > 0 {
			return math.Log(factor * float64(counts[i]) / float64(counts[i-1]))
		}

		factor *= alpha
	}

	return UnknownWordScore
}
