package lm

// ScorerNext represents the entity that responses for scoring the word using the parent context
type ScorerNext interface {
	// ScoreNext calculates the score for the given nGram built on the parent context
	ScoreNext(nGram WordID) float64
}

type scorerNext struct {
	contextCounts []WordCount
	context       ContextOffset
	nGramVector   NGramVector
}

func (s *scorerNext) ScoreNext(nGram WordID) float64 {
	count, _ := s.nGramVector.GetCount(nGram, s.context)

	if count == 0 {
		return UnknownWordScore
	}

	return calcScore(append(s.contextCounts, count))
}
