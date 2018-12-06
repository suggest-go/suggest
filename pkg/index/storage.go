package index

// IndicesReader represents entity for loading InvertedIndexIndices from storage
type IndicesReader interface {
	// Load loads inverted index indices structure from disk
	Load() (InvertedIndexIndices, error)
}

// IndicesWriter represents entity for saving Index structure in storage
type IndicesWriter interface {
	// Save tries to save indices on disc, return non nil error on failure
	Save(indices Indices) error
}
