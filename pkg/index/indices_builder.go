package index

import "github.com/alldroll/suggest/pkg/dictionary"

// IndicesBuilder is the entity that is resposible for building Indices from the given dictionary
type IndicesBuilder interface {
	// Build builds Indices from the given dictionary
	Build(dictionary dictionary.Dictionary) (Indices, error)
}

// NewIndicesBuilder returns new instance of IndicesBuilder
func NewIndicesBuilder(
	nGramSize int,
	generator Generator,
	cleaner Cleaner,
) IndicesBuilder {
	return &indicesBuilder{
		nGramSize: nGramSize,
		generator: generator,
		cleaner:   cleaner,
	}
}

// indicesBuilder implements IndicesBuilder interface
type indicesBuilder struct {
	nGramSize int
	generator Generator
	cleaner   Cleaner
}

// Build builds Indices from the given dictionary
func (ix *indicesBuilder) Build(dict dictionary.Dictionary) (Indices, error) {
	indices := make(Indices, 1)
	indices[0] = make(Index)

	err := dict.Iterate(func(key dictionary.Key, word dictionary.Value) {
		if len(word) >= ix.nGramSize {
			prepared := ix.cleaner.CleanAndWrap(word)
			set := ix.generator.Generate(prepared)
			cardinality := len(set)

			if len(indices) <= cardinality {
				tmp := make(Indices, cardinality+1, cardinality*2)
				copy(tmp, indices)
				indices = tmp
			}

			index := indices[cardinality]
			if index == nil {
				index = make(Index)
				indices[cardinality] = index
			}

			for _, term := range set {
				index[term] = append(index[term], key)
				indices[0][term] = append(indices[0][term], key)
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return indices, nil
}
