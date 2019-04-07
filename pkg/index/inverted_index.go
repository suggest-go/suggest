package index

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/compression"
)

type (
	// DocumentID is a unique indentificator of a indexed document
	DocumentID = uint32
	// Term represents an independent search element in a search document
	Term = string
	// Position (posting) is a list item of PostingList
	Position = DocumentID
	// PostingList is a list of "documents", that contains the specific term
	PostingList = []Position
	// Index is a low level data structure for storing a map of posting lists
	Index = map[Term]PostingList
	// Indices is a list of Indexes grouped by a length of a document's nGram set
	Indices = []Index
)

// InvertedIndex is an index data structure that contains list of
// references to documents for each term
type InvertedIndex interface {
	// Get returns corresponding posting list for given term
	Get(term Term) (PostingList, error)
	// Has checks is there is given term in inverted index
	Has(term Term) bool
}

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

// invertedIndexStructure is a first part of inverted index. It is stored in memory,
// with pointers to each posting list, which is stored on disk
type invertedIndexStructure map[Term]struct {
	// size - byte length of position list for given term
	size uint32
	// position - is file position, where positing list is stored
	position uint32
}

// NewInvertedIndex returns new instance of InvertedIndex that is stored on disc
func NewInvertedIndex(reader Input, decoder compression.Decoder, m invertedIndexStructure) InvertedIndex {
	return &invertedIndexImpl{
		reader:  reader,
		decoder: decoder,
		m:       m,
	}
}

// invertedIndexImpl implements InvertedIndex interface
type invertedIndexImpl struct {
	reader  Input
	decoder compression.Decoder
	m       invertedIndexStructure
}

// Get returns corresponding posting list for given term
func (i *invertedIndexImpl) Get(term Term) (PostingList, error) {
	s, ok := i.m[term]
	if !ok {
		return nil, nil
	}

	buf, err := i.reader.Data()

	if err != nil {
		return nil, fmt.Errorf("Failed to read data from index.Input: %v", err)
	}

	return i.decoder.Decode(buf[s.position : s.position+s.size]), nil
}

// Has checks is there is given term in inverted index
func (i *invertedIndexImpl) Has(term Term) bool {
	_, ok := i.m[term]
	return ok
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
