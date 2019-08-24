package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"

	"github.com/alldroll/suggest/pkg/index"
)

// Builder is the entity that is responsible for tuning and creating a NGramIndex
type Builder interface {
	// Build configures and returns a new instance of NGramIndex
	Build() (NGramIndex, error)
}

// builderImpl implements Builder interface
type builderImpl struct {
	indexReader *index.Reader
	description IndexDescription
}

// NewRAMBuilder creates a search index by using the given dictionary and the index description
// in a RAMDriver directory
func NewRAMBuilder(dict dictionary.Dictionary, description IndexDescription) (Builder, error) {
	directory := store.NewRAMDirectory()

	if err := Index(directory, dict, description.GetWriterConfig(), description.GetIndexTokenizer()); err != nil {
		return nil, fmt.Errorf("failed to create a ram search index: %v", err)
	}

	return NewBuilder(directory, description)
}

// NewFSBuilder works with already indexed data
func NewFSBuilder(description IndexDescription) (Builder, error) {
	directory, err := store.NewFSDirectory(description.GetIndexPath())

	if err != nil {
		return nil, fmt.Errorf("failed to create a fs directory: %v", err)
	}

	return NewBuilder(directory, description)
}

// NewBuilder works with already indexed data
func NewBuilder(directory store.Directory, description IndexDescription) (Builder, error) {
	return &builderImpl{
		indexReader: index.NewIndexReader(
			directory,
			description.GetWriterConfig(),
		),
		description: description,
	}, nil
}

// Build configures and returns a new instance of NGramIndex
func (b *builderImpl) Build() (NGramIndex, error) {
	invertedIndices, err := b.indexReader.Read()

	if err != nil {
		return nil, fmt.Errorf("failed to build NGramIndex: %v", err)
	}

	suggester := NewSuggester(
		invertedIndices,
		index.NewSearcher(merger.CPMerge()),
		NewSuggestTokenizer(b.description),
	)

	autocomplete := NewAutocomplete(
		invertedIndices.GetWholeIndex(),
		index.NewSearcher(merger.CPMerge()),
		NewAutocompleteTokenizer(b.description),
	)

	return NewNGramIndex(
		suggester,
		autocomplete,
	), nil
}
