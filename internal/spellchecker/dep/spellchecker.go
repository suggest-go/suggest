package dep

import (
	"fmt"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/spellchecker"
	"github.com/suggest-go/suggest/pkg/store"
	"github.com/suggest-go/suggest/pkg/suggest"
)

// BuildSpellChecker builds spellchecker for the provided config and indexDescription
func BuildSpellChecker(config *lm.Config, indexDescription suggest.IndexDescription) (*spellchecker.SpellChecker, error) {
	directory, err := store.NewFSDirectory(config.GetOutputPath())

	if err != nil {
		return nil, fmt.Errorf("failed to create a fs directory: %v", err)
	}

	languageModel, err := lm.RetrieveLMFromBinary(directory, config)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a lm model from binary format: %v", err)
	}

	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		return nil, fmt.Errorf("failed to open a cdb dictionary: %v", err)
	}

	// create runtime search index builder
	builder, err := suggest.NewRAMBuilder(dict, indexDescription)

	if err != nil {
		return nil, fmt.Errorf("failed to create a ngram index: %v", err)
	}

	index, err := builder.Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build a ngram index: %v", err)
	}

	return spellchecker.New(
		index,
		languageModel,
		lm.NewTokenizer(config.GetWordsAlphabet()),
		dict,
	), nil
}
