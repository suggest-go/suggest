package index

// Indices is a list of Indexes grouped by a length of a document's nGram set
type Indices = []Index

// InvertedIndexIndices is a array of InvertedIndex, where index - ngrams cardinality of containing documents
// 0 index - inverted index that contains all documents (without ngrams' cardinality separation)
type InvertedIndexIndices interface {
	// GetWholeIndex returns whole InvertedIndex (without ngram's cardinality separation)
	GetWholeIndex() InvertedIndex
	// Get returns InvertedIndex of term with given index.
	// Index here represents document ngrams cardinality
	Get(index int) InvertedIndex
	// Size returns number of InvertedIndex
	Size() int
}

// NewInvertedIndexIndices returns new instance of InvertedIndexIndices
func NewInvertedIndexIndices(indices []InvertedIndex) InvertedIndexIndices {
	return &invertedIndexIndicesImpl{indices}
}

// invertedIndexIndicesImpl implements InvertedIndexIndices interface
type invertedIndexIndicesImpl struct {
	indices []InvertedIndex
}

// GetWholeIndex returns whole InvertedIndex (without ngram's cardinality separation)
func (i *invertedIndexIndicesImpl) GetWholeIndex() InvertedIndex {
	return i.indices[0]
}

// Get returns InvertedIndex of term with given index.
func (i *invertedIndexIndicesImpl) Get(index int) InvertedIndex {
	if index >= 0 && index < len(i.indices) {
		return i.indices[index]
	}

	return nil
}

// Size returns number of InvertedIndex
func (i *invertedIndexIndicesImpl) Size() int {
	return len(i.indices)
}
