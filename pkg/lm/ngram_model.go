package lm

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
)

// NGramModel is an entity that responses for scoring the given nGrams
type NGramModel interface {
	// Score returns a lm value of the given sequence of WordID
	Score(nGrams []WordID) float64
	// Next returns a list of WordID which follow after the given sequence of nGrams
	Next(nGrams []WordID) (ScorerNext, error)
}

const (
	// UnknownWordScore is the score for unknown phrases
	UnknownWordScore = -100.0
	alpha            = 0.4
	version          = "0.0.1"
)

// nGramModel implements NGramModel Stupid backoff
type nGramModel struct {
	indices    []NGramVector
	nGramOrder uint8
}

// NewNGramModel creates a new instance of NGramModel instance
func NewNGramModel(indices []NGramVector) NGramModel {
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
		nGramVector: subVector,
		context: parent,
	}, nil
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (m *nGramModel) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(version); err != nil {
		return nil, err
	}

	if err := encoder.Encode(m.nGramOrder); err != nil {
		return nil, err
	}

	for _, vector := range m.indices {
		err := encoder.Encode(&vector)

		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the binary form
func (m *nGramModel) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	binaryVersion := "NONE"

	if err := decoder.Decode(&binaryVersion); err != nil {
		return err
	}

	if binaryVersion != version {
		return fmt.Errorf("version mismatch, expected: %s, got %s", version, binaryVersion)
	}

	if err := decoder.Decode(&m.nGramOrder); err != nil {
		return err
	}

	m.indices = make([]NGramVector, int(m.nGramOrder))

	for i := 0; i < len(m.indices); i++ {
		err := decoder.Decode(&m.indices[i])

		if err != nil {
			return err
		}
	}

	return nil
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

func init() {
	gob.Register(&nGramModel{})
}
