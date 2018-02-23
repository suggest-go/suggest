package index

// IndicesReader represents entity for loading InvertedIndex from storage
type IndexReader interface {
	// Load loads inverted index structure from disk
	Load() (InvertedIndex, error)
}

// IndicesWriter represents entity for saving Index structure in storage
type IndexWriter interface {
	// Save tries to save index on disc, return non nil error on failure
	Save(index Index) error
}

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
