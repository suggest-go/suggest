package suggest

type EditDistance interface {
	Calc(profileA, profileB *WordProfile) float64
}

type LevenshteinDistance struct{}

func (self *LevenshteinDistance) calc(a, b string) float64 {
	r1, r2 := []rune(a), []rune(b)
	aLen, bLen := len(r1), len(r2)
	if aLen == 0 {
		return float64(bLen)
	}

	if bLen == 0 {
		return float64(aLen)
	}

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

func (self *LevenshteinDistance) Calc(profileA, profileB *WordProfile) float64 {
	return self.calc(profileA.GetWord(), profileB.GetWord())
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
func (self *NGramDistance) Calc(profileA, profileB *WordProfile) float64 {
	distance := 0.0
	//for _, key := range profileA.ngrams {
	//freqA, freqB := profileA.frequencies[key], 0
	//if val, ok := profileB.frequencies[key]; ok {
	//freqB = val
	//}

	//distance += math.Abs(float64(freqA - freqB))
	//}

	//for _, key := range profileB.ngrams {
	//if _, ok := profileA.frequencies[key]; !ok {
	//distance += float64(profileB.frequencies[key])
	//}
	//}

	return distance
}

type JaccardDistance struct {
	k int
}

func CreateJaccardDistance(k int) *JaccardDistance {
	return &JaccardDistance{k}
}

// Jaccard distance = 1 - J(A, B) = 1 - |intersection| / |union|
func (self *JaccardDistance) Calc(profileA, profileB *WordProfile) float64 {
	if profileA.GetWord() == profileB.GetWord() {
		return 0.0
	}

	ngramsA, ngramsB := profileA.GetNGrams(), profileB.GetNGrams()
	lenA, lenB := len(ngramsA), len(ngramsB)
	i, j := 0, 0
	inter := 0.0
	for i < lenA && j < lenB {
		aVal, bVal := ngramsA[i].GetValue(), ngramsB[j].GetValue()
		if aVal < bVal {
			i++
		} else if bVal < aVal {
			j++
		} else {
			inter++
			i++
			j++
		}
	}

	return 1.0 - inter/(float64(lenA+lenB)-inter) //union = |a|+|b|-intersection
}
