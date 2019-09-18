// Package spellchecker provides spellcheck functionality
package spellchecker

import (
	"github.com/alldroll/suggest/pkg/analysis"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/lm"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"
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
	collectorManager, err := s.createCollectorManager(seq, topK)

	if err != nil {
		return nil, err
	}

	candidates, err := s.index.Autocomplete(word, collectorManager)

	if err != nil {
		return nil, err
	}

	if len(candidates) < topK {
		config, err := suggest.NewSearchConfig(
			word,
			topK-len(candidates),
			metric.CosineMetric(),
			similarity,
		)

		if err != nil {
			return nil, err
		}

		fuzzyCandidates, err := s.index.Suggest(config)

		if err != nil {
			return nil, err
		}

		candidates = merge(candidates, fuzzyCandidates)
	}

	if len(seq) > 0 {
		sortCandidates(collectorManager.scorer, candidates)
	}

	result := make([]string, 0, len(candidates))

	for _, c := range candidates {
		val, err := s.dict.Get(c.Key)

		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

// createScorer creates scorer for the given sentence
func (s *SpellChecker) createCollectorManager(seq []string, topK int) (*lmCollectorManager, error) {
	seqIds, err := lm.MapIntoListOfWordIDs(s.model, seq)

	if err != nil {
		return nil, err
	}

	next := []lm.WordID{}

	if len(seqIds) > 0 {
		next, err = s.model.Next(seqIds)

		if err != nil {
			return nil, err
		}
	}

	return &lmCollectorManager{
		topK: topK,
		scorer: &lmScorer{
			model:    s.model,
			sentence: seqIds,
		},
		next: next,
	}, nil
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
