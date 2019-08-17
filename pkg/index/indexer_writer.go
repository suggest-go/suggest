package index

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
)

// IndexVersion tells that the inverted index structure has the provided below version
const IndexVersion = "v5"

// Writer creates and maintains an inverted index
type Writer struct {
	directory store.Directory
	config    WriterConfig
	encoder   compression.Encoder
	indices   Indices
}

// WriterConfig stores a set of file paths that are required
// for creating search index
type WriterConfig struct {
	HeaderFileName       string
	DocumentListFileName string
}

// NewIndexWriter returns new instance of a index writer
func NewIndexWriter(
	directory store.Directory,
	config WriterConfig,
	encoder compression.Encoder,
) *Writer {
	return &Writer{
		directory: directory,
		config:    config,
		encoder:   encoder,
		indices:   Indices{},
	}
}

var (
	// ErrPostingListShouldBeNotNil occurs when was an attempt to persist nil Posting List
	ErrPostingListShouldBeNotNil = errors.New("PostingList should be not nil")
)

// header struct that store terms descriptions and indices count
type header struct {
	Version string
	Indices uint32
	Terms   []termDescription
}

// termDescription stores term, indice, postingList size and postingList file position
type termDescription struct {
	Term                 Term
	Indice               uint32
	PostingListBytesSize uint32
	PostingListPosition  uint32
	PostingListLen       uint32
}

// AddDocument adds a new documents with the given fields
func (iw *Writer) AddDocument(id DocumentID, term []Term) error {
	cardinality := len(term)

	if len(iw.indices) <= cardinality {
		tmp := make(Indices, cardinality+1, cardinality*2)
		copy(tmp, iw.indices)
		iw.indices = tmp
	}

	index := iw.indices[cardinality]
	if index == nil {
		index = make(Index)
		iw.indices[cardinality] = index
	}

	if iw.indices[0] == nil {
		iw.indices[0] = make(Index)
	}

	for _, term := range term {
		index[term] = append(index[term], id)
		iw.indices[0][term] = append(iw.indices[0][term], id)
	}

	return nil
}

// Commit commits all added documents to the index storage
func (iw *Writer) Commit() error {
	documentWriter, err := iw.directory.CreateOutput(iw.config.DocumentListFileName)

	if err != nil {
		return fmt.Errorf("Failed to create document list: %v", err)
	}

	// mapValueOffset stores current posting list offset
	mapValueOffset := int64(0)

	// header struct that should be loaded on Load
	header := header{
		Version: IndexVersion,
		Terms:   []termDescription{},
		Indices: uint32(len(iw.indices)),
	}

	for indice, index := range iw.indices {
		if index == nil {
			continue
		}

		for term, postingList := range index {
			// there is not possible, we should throw the error
			if postingList == nil {
				return ErrPostingListShouldBeNotNil
			}

			// Encode the given posting list into a byte slice
			n, err := iw.encoder.Encode(postingList, documentWriter)

			if err != nil {
				return err
			}

			header.Terms = append(header.Terms, termDescription{
				Term:                 term,
				Indice:               uint32(indice),
				PostingListBytesSize: uint32(n),
				PostingListPosition:  uint32(mapValueOffset),
				PostingListLen:       uint32(len(postingList)),
			})

			mapValueOffset += int64(n)
		}
	}

	if err = iw.writeHeader(header); err != nil {
		return err
	}

	if err = documentWriter.Close(); err != nil {
		return fmt.Errorf("Failed to close document list: %v", err)
	}

	return nil
}

// writeHeader writes and persists index header
func (iw *Writer) writeHeader(header header) error {
	headerWriter, err := iw.directory.CreateOutput(iw.config.HeaderFileName)

	if err != nil {
		return fmt.Errorf("Failed to create header: %v", err)
	}

	encoder := gob.NewEncoder(headerWriter)
	err = encoder.Encode(header)

	if err != nil {
		return fmt.Errorf("Failed to encode header: %v", err)
	}

	if err = headerWriter.Close(); err != nil {
		return fmt.Errorf("Failed to close header file: %v", err)
	}

	return nil
}
