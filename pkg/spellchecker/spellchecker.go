package spellchecker

import (
	"github.com/alldroll/suggest/pkg/dictionary"
	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"
)

// SpellChecker describe me!
type SpellChecker struct {
	index     suggest.NGramIndex
	model     lm.LanguageModel
	tokenizer lm.Tokenizer
	dict      dictionary.Dictionary
}

// New creates a new instance of spellchecker
func New(
	index suggest.NGramIndex,
	model lm.LanguageModel,
	tokenizer lm.Tokenizer,
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
func (s *SpellChecker) Predict(query string, topK int) ([]string, error) {
	tokens := s.tokenizer.Tokenize(query)

	if len(tokens) == 0 {
		return []string{}, nil
	}

	word, seq := tokens[len(tokens)-1], tokens[:len(tokens)-1]
	seqIds, err := lm.MapIntoListOfWordIDs(s.model, seq)

	if err != nil {
		return nil, err
	}

	scorer := &lmScorer{
		model:    s.model,
		sentence: seqIds,
	}

	candidates, err := s.index.AutoComplete(word, topK, scorer)

	if err != nil {
		return nil, err
	}

	if len(candidates) < topK {
		config, err := suggest.NewSearchConfig(
			word,
			topK-len(candidates),
			metric.CosineMetric(),
			0.7,
		)

		if err != nil {
			return nil, err
		}

		fuzzyCandidates, err := s.index.Suggest(config)
		candidates = append(candidates, fuzzyCandidates...)
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
