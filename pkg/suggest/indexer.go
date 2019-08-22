package suggest

import (
	"fmt"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/store"
)

// Index builds a search index by using the given config and the dictionary
// and persists it the directory
func Index(directory store.Directory, dict dictionary.Dictionary, config IndexDescription) error {
	alphabet := config.CreateAlphabet()
	cleaner, err := index.NewCleaner(alphabet.Chars(), config.Pad, config.Wrap)

	if err != nil {
		return err
	}

	encoder, err := index.NewEncoder()

	if err != nil {
		return fmt.Errorf("failed to create Encoder: %v", err)
	}

	generator := index.NewGenerator(config.NGramSize)
	indexWriter := index.NewIndexWriter(
		directory,
		config.CreateWriterConfig(),
		encoder,
	)

	if err = index.BuildIndex(dict, indexWriter, generator, cleaner); err != nil {
		return err
	}

	return nil
}
