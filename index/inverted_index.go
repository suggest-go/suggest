package index

import (
	"github.com/alldroll/suggest/compression"
	"golang.org/x/exp/mmap"
	"io"
	"runtime"
)

type (
	Term = uint32
	// Position (posting) is list item of PostingList
	Position = DocumentID
	// PostingList is a list of "documents", which contains specific term
	PostingList = []Position
)

// InvertedIndex is an index data structure that contains list of
// references to documents for each term
type InvertedIndex interface {
	// Get returns corresponding posting list for given term
	Get(term Term) PostingList
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

// InvertedIndexIndicesBuilder
type InvertedIndexIndicesBuilder interface {
	// Build
	Build() InvertedIndexIndices
}

// NewInMemoryInvertedIndex returns new instance of InvertedIndex that is stored in memory
func NewInMemoryInvertedIndex(index Index) InvertedIndex {
	return &invertedIndexInMemoryImpl{index}
}

// invertedIndexInMemoryImpl is in memory inverted index implementation
type invertedIndexInMemoryImpl struct {
	table map[Term]PostingList
}

// Get returns corresponding posting list for given term
func (i *invertedIndexInMemoryImpl) Get(term Term) PostingList {
	return i.table[term]
}

// Has checks is there is given term in inverted index
func (i *invertedIndexInMemoryImpl) Has(term Term) bool {
	_, ok := i.table[term]
	return ok
}

// invertedIndexStructure is a first part of inverted index. It is stored in memory,
// with pointers to each posting list, which is stored on disk
type invertedIndexStructure map[Term]struct {
	// size - byte length of position list for given term
	size uint32
	// position - is file position, where positing list is stored
	position uint32
}

// NewOnDiscInvertedIndex returns new instance of InvertedIndex that is stored on disc
func NewOnDiscInvertedIndex(reader io.ReaderAt, decoder compression.Decoder, m invertedIndexStructure) InvertedIndex {
	return &onDiscInvertedIndex{
		reader:  reader,
		decoder: decoder,
		m:       m,
	}
}

// onDiscInvertedIndex is on disk inverted index implementation
type onDiscInvertedIndex struct {
	reader  io.ReaderAt
	decoder compression.Decoder
	m       invertedIndexStructure
}

// Get returns corresponding posting list for given term
func (i *onDiscInvertedIndex) Get(term Term) PostingList {
	s, ok := i.m[term]
	if !ok {
		return nil
	}

	buf := make([]byte, s.size)
	i.reader.ReadAt(buf, int64(s.position))

	return i.decoder.Decode(buf)
}

// Has checks is there is given term in inverted index
func (i *onDiscInvertedIndex) Has(term Term) bool {
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

// NewInMemoryInvertedIndexIndicesBuilder returns new instance of InvertedIndexIndicesBuilder
// that builds "In memory" InvertedIndexIndices
func NewInMemoryInvertedIndexIndicesBuilder(indices Indices) InvertedIndexIndicesBuilder {
	return &invertedIndexIndicesBuilderInMemoryImpl{
		indices: indices,
	}
}

// invertedIndexIndicesBuilderInMemoryImpl builds "In memory" InvertedIndexIndices
type invertedIndexIndicesBuilderInMemoryImpl struct {
	indices Indices
}

// Build creates new "In memory" InvertedIndexIndices
func (b *invertedIndexIndicesBuilderInMemoryImpl) Build() InvertedIndexIndices {
	invertedIndexIndices := make([]InvertedIndex, len(b.indices))

	for i, index := range b.indices {
		invertedIndexIndices[i] = NewInMemoryInvertedIndex(index)
	}

	return NewInvertedIndexIndices(invertedIndexIndices)
}

// NewOnDiscInvertedIndexIndicesBuilder returns new instance of InvertedIndexIndicesBuilder
// that builds "On disc" InvertedIndexIndices
func NewOnDiscInvertedIndexIndicesBuilder(headerPath, documentListPath string) InvertedIndexIndicesBuilder {
	return &invertedIndexIndicesBuilderOnDiscImpl{
		headerPath:       headerPath,
		documentListPath: documentListPath,
	}
}

// invertedIndexIndicesBuilderOnDiscImpl builds "On disc" InvertedIndexIndices
type invertedIndexIndicesBuilderOnDiscImpl struct {
	// headerPath is path to file, that contains InvertedIndexIndices structure
	headerPath string
	// documentListPath is path to file, that contains PostingList
	documentListPath string
}

// Build creates new "On disc" InvertedIndexIndices
func (b *invertedIndexIndicesBuilderOnDiscImpl) Build() InvertedIndexIndices {
	header, err := mmap.Open(b.headerPath)
	if err != nil {
		panic(err)
	}

	docList, err := mmap.Open(b.documentListPath)
	if err != nil {
		header.Close()
		panic(err)
	}

	reader := NewOnDiscIndicesReader(compression.VBDecoder(), header, docList, 0)
	indices, err := reader.Load()
	if err != nil {
		header.Close()
		docList.Close()
		panic(err)
	}

	runtime.SetFinalizer(indices, func(impl *invertedIndexIndicesImpl) {
		header.Close()
		docList.Close()
	})

	return indices
}
