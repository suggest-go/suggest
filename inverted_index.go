package suggest

import (
	"golang.org/x/exp/mmap"
	"io"
	"runtime"
)

type Term int32
type Position uint32
type PostingList []Position

// InvertedIndex
type InvertedIndex interface {
	// Get returns corresponding posting list for given term
	Get(term Term) PostingList
	// Has checks is there is given term in inverted index
	Has(term Term) bool
}

// InvertedIndexIndices
type InvertedIndexIndices interface {
	// Get
	Get(index int) InvertedIndex
	// Size
	Size() int
}

// InvertedIndexIndicesBuilder
type InvertedIndexIndicesBuilder interface {
	// Build
	Build() InvertedIndexIndices
}

// NewInMemoryInvertedIndex
func NewInMemoryInvertedIndex(table map[Term]PostingList) InvertedIndex {
	return &invertedIndexInMemoryImpl{table}
}

type invertedIndexStructure map[Term]struct {
	size     uint32
	position uint32
}

func NewOnDiscInvertedIndex(reader io.ReaderAt, decoder Decoder, m invertedIndexStructure) InvertedIndex {
	return &onDiscInvertedIndex{
		reader:  reader,
		decoder: decoder,
		m:       m,
	}
}

// invertedIndexInMemoryImpl
type invertedIndexInMemoryImpl struct {
	table map[Term]PostingList
}

// Get
func (i *invertedIndexInMemoryImpl) Get(term Term) PostingList {
	return i.table[term]
}

func (i *invertedIndexInMemoryImpl) Has(term Term) bool {
	_, ok := i.table[term]
	return ok
}

// onDiscInvertedIndex
type onDiscInvertedIndex struct {
	reader  io.ReaderAt
	decoder Decoder
	m       invertedIndexStructure
}

func (i *onDiscInvertedIndex) Get(term Term) PostingList {
	s, ok := i.m[term]
	if !ok {
		return nil
	}

	buf := make([]byte, s.size)
	i.reader.ReadAt(buf, int64(s.position))

	return i.decoder.Decode(buf)
}

func (i *onDiscInvertedIndex) Has(term Term) bool {
	_, ok := i.m[term]
	return ok
}

// NewInvertedIndexIndices
func NewInvertedIndexIndices(indices []InvertedIndex) InvertedIndexIndices {
	return &invertedIndexIndicesImpl{indices}
}

// invertedIndexIndicesImpl
type invertedIndexIndicesImpl struct {
	indices []InvertedIndex
}

// Get
func (i *invertedIndexIndicesImpl) Get(index int) InvertedIndex {
	if index >= 0 && index < len(i.indices) {
		return i.indices[index]
	}

	return nil
}

// Size
func (i *invertedIndexIndicesImpl) Size() int {
	return len(i.indices)
}

// NewInMemoryInvertedIndexIndicesBuilder
func NewInMemoryInvertedIndexIndicesBuilder(indices Index) InvertedIndexIndicesBuilder {
	return &invertedIndexIndicesBuilderInMemoryImpl{
		indices: indices,
	}
}

// invertedIndexIndicesBuilderInMemoryImpl
type invertedIndexIndicesBuilderInMemoryImpl struct {
	indices Index
}

// Build
func (b *invertedIndexIndicesBuilderInMemoryImpl) Build() InvertedIndexIndices {
	invertedIndexIndices := make([]InvertedIndex, len(b.indices))

	for i, table := range b.indices {
		invertedIndexIndices[i] = NewInMemoryInvertedIndex(table)
	}

	return NewInvertedIndexIndices(invertedIndexIndices)
}

// NewOnDiscInvertedIndexIndicesBuilder
func NewOnDiscInvertedIndexIndicesBuilder(headerPath, documentListPath string) InvertedIndexIndicesBuilder {
	return &invertedIndexIndicesBuilderOnDiscImpl{
		headerPath:       headerPath,
		documentListPath: documentListPath,
	}
}

// invertedIndexIndicesBuilderOnDiscImpl
type invertedIndexIndicesBuilderOnDiscImpl struct {
	headerPath, documentListPath string
}

// Build (monkey code, fix me)
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

	reader := NewOnDiscInvertedIndexReader(VBDecoder(), header, docList, 0)
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
