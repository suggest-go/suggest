package suggest

import (
	"github.com/alldroll/cdb"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

type PostingList []int

// InvertedIndex
type InvertedIndex interface {
	// Get
	Get(term int) PostingList
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
func NewInMemoryInvertedIndex(table map[int]PostingList) InvertedIndex {
	return &invertedIndexInMemoryImpl{table}
}

// invertedIndexInMemoryImpl
type invertedIndexInMemoryImpl struct {
	table map[int]PostingList
}

// Get
func (i *invertedIndexInMemoryImpl) Get(term int) PostingList {
	return i.table[term]
}


// NewCdbInvertedIndex
func NewCdbInvertedIndex(reader cdb.Reader) InvertedIndex {
	return &invertedIndexCDBImpl{reader}
}

//
type invertedIndexCDBImpl struct {
	reader cdb.Reader
}

// Get
func (i *invertedIndexCDBImpl) Get(term int) PostingList {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(term))

	d, err := i.reader.Get(b)
	if err != nil {
		// TODO handle me
		panic(err)
	}

	list := make([]int, len(d) / 4)
	for i := range list {
		list[i] = int(binary.LittleEndian.Uint32(d[4*i:]))
	}

	return list
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
	return &invertedIndexIndicesBuilderInMemoryImpl{indices}
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

// NewCDBInvertedIndexIndicesBuilder
func NewCDBInvertedIndexIndicesBuilder(dbDir string) InvertedIndexIndicesBuilder {
	return &invertedIndexIndicesBuilderCDBImpl{dbDir}
}

// invertedIndexIndicesBuilderCDBImpl
type invertedIndexIndicesBuilderCDBImpl struct {
	dbDir string
}

// Build (monkey code, fix me)
func (b *invertedIndexIndicesBuilderCDBImpl) Build() InvertedIndexIndices {
	cdbHandle := cdb.New()
	indices := make([]InvertedIndex, 0)

	err := filepath.Walk(b.dbDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || err != nil {
			return err
		}

		baseName := info.Name()
		indexStr := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}

		reader, err := cdbHandle.GetReader(f)
		if err != nil {
			return err
		}

		if len(indices) <= index {
			tmp := make([]InvertedIndex, index+1, index*2)
			copy(tmp, indices)
			indices = tmp
		}

		indices[index] = NewCdbInvertedIndex(reader)
		return err
	})

	if err != nil {
		panic(err)
	}

	return NewInvertedIndexIndices(indices)
}
