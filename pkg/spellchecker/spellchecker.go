// Package spellchecker provides spellcheck functionality
package spellchecker

// TODO add tests!!

import (
	"sort"

	"github.com/suggest-go/suggest/pkg/analysis"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
)

// SpellChecker describe me!
type SpellChecker struct {
	index     suggest.NGramIndex
	model     lm.LanguageModel
	tokenizer analysis.Tokenizer
	dict      dictionary.Dictionary
}

// New creates a new instance of spellchecker
func New(
	index suggest.NGramIndex,
	model lm.LanguageModel,
	tokenizer analysis.Tokenizer,
	dict dictionary.Dictionary,
) *SpellChecker {
	return &SpellChecker{
		index:     index,
		model:     model,
		tokenizer: tokenizer,
		dict:      dict,
	}
}

// Predict predicts the next word of the sentence
func (s *SpellChecker) Predict(query string, topK int, similarity float64) ([]string, error) {
	tokens := s.tokenizer.Tokenize(query)

	if len(tokens) == 0 {
		return []string{}, nil
	}

	word, seq := tokens[len(tokens)-1], tokens[:len(tokens)-1]
	scorerNext, err := scorerNext(s.model, seq)

	if err != nil {
		return nil, err
	}

	queueFactory := func() suggest.TopKQueue {
		return suggest.NewTopKQueue(topK)
	}

	candidates, err := s.index.Autocomplete(word, func() suggest.CollectorManager {
		return newCollectorManager(newScorer(scorerNext), queueFactory)
	})

	if err != nil {
		return nil, err
	}

	if len(candidates) < topK {
		fuzzyCandidates, err := s.index.Suggest(
			word,
			similarity,
			metric.CosineMetric(),
			func() suggest.CollectorManager {
				return suggest.NewFuzzyCollectorManager(queueFactory)
			},
		)

		if err != nil {
			return nil, err
		}

		candidates = merge(candidates, fuzzyCandidates)
	}

	if scorerNext != nil {
		sortCandidates(scorerNext, candidates)
	}

	if topK < len(candidates) {
		candidates = candidates[:topK+1]
	}

	return retrieveValues(s.dict, candidates)
}

// scorerNext creates lm.ScorerNext for the provided sentence
func scorerNext(model lm.LanguageModel, seq lm.Sentence) (next lm.ScorerNext, err error) {
	seqIds, err := lm.MapIntoListOfWordIDs(model, seq)

	if err != nil {
		return nil, err
	}

	if len(seqIds) > 0 {
		next, err = model.Next(seqIds)
	}

	return
}

// retrieveValues fetches the corresponding values from the dictionary for the provided candidates.
func retrieveValues(dict dictionary.Dictionary, candidates []suggest.Candidate) ([]string, error) {
	result := make([]string, 0, len(candidates))

	for _, c := range candidates {
		val, err := dict.Get(c.Key)

		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

// sortCandidates performs sort of the given candidates using lm
func sortCandidates(scorer lm.ScorerNext, candidates []suggest.Candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		return scorer.ScoreNext(candidates[i].Key) > scorer.ScoreNext(candidates[j].Key)
	})
}

// merge merges the 2 candidates sets into one without duplication
func merge(a, b []suggest.Candidate) []suggest.Candidate {
	for _, y := range b {
		unique := true

		for _, x := range a {
			if x.Key == y.Key {
				unique = false
				break
			}
		}

		if unique {
			a = append(a, y)
		}
	}

	return a
}
