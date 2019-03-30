package lm

type (
	// NGrams is the result of splitting the given sequence of words into nGrams
	NGrams = [][]WordID
)

// SplitIntoNGrams splits the given sequence of WordID into a set of nGrams
func SplitIntoNGrams(sequence []WordID, nGramOrder uint8) NGrams {
	k := int(nGramOrder)

	if len(sequence) < k {
		return NGrams{}
	}

	nGrams := make(NGrams, 0, len(sequence)-k+1)

	for i := 0; i <= len(sequence)-k; i++ {
		nGrams = append(nGrams, sequence[i:i+k])
	}

	return nGrams
}
