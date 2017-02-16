package suggest

import (
	"errors"
	"math"
)

const (
	LEVENSHTEIN = iota
	NGRAM
	JACCARD
)

var MetricName = map[int]string{
	LEVENSHTEIN: `levenshtein`,
	NGRAM:       `ngram`,
	JACCARD:     `jaccard`,
}

type EditDistance interface {
	Calc(a, b string) float64
	/* monkeycode */
	CalcWithProfiles(a, b string, profileA, profileB *profile) float64
}

func GetEditDistance(t int, k int) (EditDistance, error) {
	switch t {
	case LEVENSHTEIN:
		return &LevenshteinDistance{}, nil

	case NGRAM:
		return &NGramDistance{k}, nil

	case JACCARD:
		return &JaccardDistance{k}, nil

	default:
		return nil, errors.New("Invalid metric type")
	}
}

type LevenshteinDistance struct{}

func (self *LevenshteinDistance) Calc(a, b string) float64 {
	aLen, bLen := len(a), len(b)
	if aLen == 0 {
		return float64(bLen)
	}

	if bLen == 0 {
		return float64(aLen)
	}

	r1, r2 := []rune(a), []rune(b)
	column := make([]int, aLen+1)

	for i := 1; i < aLen+1; i++ {
		column[i] = i
	}

	for j := 1; j < bLen+1; j++ {
		column[0] = j
		prev := j - 1
		for i := 1; i < aLen+1; i++ {
			tmp := column[i]
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			column[i] = min3(
				column[i]+1,
				column[i-1]+1,
				prev+cost,
			)
			prev = tmp
		}
	}

	return float64(column[aLen])
}

func (self *LevenshteinDistance) CalcWithProfiles(a, b string, profileA, profileB *profile) float64 {
	return self.Calc(a, b)
}

type NGramDistance struct {
	k int
}

/*
 * ngram distance between a, b defined in
 * "Approximate string-matching with q-grams and maximal matches"
 *
 * Complexity O(aLen + bLen)
 */
func (self *NGramDistance) Calc(a, b string) float64 {
	profileA, profileB := getProfile(a, self.k), getProfile(b, self.k)
	return self.CalcWithProfiles(a, b, profileA, profileB)
}

func (self *NGramDistance) CalcWithProfiles(a, b string, profileA, profileB *profile) float64 {
	distance := 0.0
	for _, key := range profileA.ngrams {
		freqA, freqB := profileA.frequencies[key], 0
		if val, ok := profileB.frequencies[key]; ok {
			freqB = val
		}

		distance += math.Abs(float64(freqA - freqB))
	}

	for _, key := range profileB.ngrams {
		if _, ok := profileA.frequencies[key]; !ok {
			distance += float64(profileB.frequencies[key])
		}
	}

	return distance
}

type JaccardDistance struct {
	k int
}

func (self *JaccardDistance) Calc(a, b string) float64 {
	if a == b {
		return 1.0
	}

	profileA, profileB := getProfile(a, self.k), getProfile(b, self.k)
	return self.CalcWithProfiles(a, b, profileA, profileB)
}

func (self *JaccardDistance) CalcWithProfiles(a, b string, profileA, profileB *profile) float64 {
	if a == b {
		return 1.0
	}

	minProfile, maxProfile := profileA, profileB
	lenA, lenB := len(profileA.frequencies), len(profileB.frequencies)
	if lenA > lenB {
		minProfile, maxProfile = maxProfile, minProfile
	}

	inter := 0.0
	for _, k := range minProfile.ngrams {
		if _, ok := maxProfile.frequencies[k]; ok {
			inter += 1
		}
	}

	return 1.0 - inter/(float64(lenA+lenB)-inter) //union = |a|+|b|-intersection
}
