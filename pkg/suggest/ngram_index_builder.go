package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/merger"
)

// Builder is the entity that is responsible for tuning and creating a NGramIndex
type Builder interface {
	// Build configures and returns a new instance of NGramIndex
	Build() (NGramIndex, error)
}

// builderImpl implements Builder interface
type builderImpl struct {
	indexReader *index.Reader
	cleaner     index.Cleaner
	generator   index.Generator
}

// NewFSBuilder works with already indexed data
func NewFSBuilder(description IndexDescription) (Builder, error) {
	alphabet := description.CreateAlphabet()
	cleaner, err := index.NewCleaner(alphabet.Chars(), description.Pad, description.Wrap)

	if err != nil {
		return nil, fmt.Errorf("Failed to create cleaner: %v", err)
	}

	generator := index.NewGenerator(description.NGramSize)
	directory, err := index.NewFSDirectory(description.OutputPath)

	if err != nil {
		return nil, fmt.Errorf("Failed to create fs directory: %v", err)
	}

	return &builderImpl{
		indexReader: index.NewIndexReader(
			directory,
			description.CreateWriterConfig(),
			compression.VBDecoder(),
		),
		generator: generator,
		cleaner:   cleaner,
	}, nil
}

// NewRAMBuilder creates a search index by using the given dictionary and the index description
// in a RAMDriver directory
func NewRAMBuilder(dict dictionary.Dictionary, description IndexDescription) (Builder, error) {
	alphabet := description.CreateAlphabet()
	cleaner, err := index.NewCleaner(alphabet.Chars(), description.Pad, description.Wrap)

	if err != nil {
		return nil, fmt.Errorf("Failed to create cleaner: %v", err)
	}

	directory := index.NewRAMDirectory()
	generator := index.NewGenerator(description.NGramSize)
	writerConfig := description.CreateWriterConfig()
	indexWriter := index.NewIndexWriter(
		directory,
		writerConfig,
		compression.VBEncoder(),
	)

	if err := index.BuildIndex(dict, indexWriter, generator, cleaner); err != nil {
		return nil, fmt.Errorf("Failed to build index in RAMDriver directory: %v", err)
	}

	return &builderImpl{
		indexReader: index.NewIndexReader(
			directory,
			writerConfig,
			compression.VBDecoder(),
		),
		generator: generator,
		cleaner:   cleaner,
	}, nil
}

// Build configures and returns a new instance of NGramIndex
func (b *builderImpl) Build() (NGramIndex, error) {
	invertedIndices, err := b.indexReader.Read()

	if err != nil {
		return nil, fmt.Errorf("Failed to build NGramIndex: %v", err)
	}

	return NewNGramIndex(
		b.cleaner,
		b.generator,
		invertedIndices,
		merger.CPMerge(),
	), nil
}
