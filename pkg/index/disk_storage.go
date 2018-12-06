package index

import (
	"bufio"
	"encoding/gob"
	"errors"
	"github.com/alldroll/suggest/pkg/compression"
	mmap "github.com/edsrzf/mmap-go"
	"os"
	"runtime"
)

var (
	PostingListShouldBeNotNilError = errors.New("PostingList should be not nil")
)

// NewOnDiscIndicesReader returns new instance of IndicesReader
func NewOnDiscIndicesReader(decoder compression.Decoder, header, documentList string) IndicesReader {
	return &onDiscIndicesReader{
		decoder:      decoder,
		header:       header,
		documentList: documentList,
	}
}

// NewOnDiscIndicesWriter returns new instance of IndicesWriter
func NewOnDiscIndicesWriter(encoder compression.Encoder, header, documentList string) IndicesWriter {
	return &onDiscIndicesWriter{
		encoder:      encoder,
		header:       header,
		documentList: documentList,
	}
}

// onDiscIndicesReader implements Reader interface
type onDiscIndicesReader struct {
	decoder              compression.Decoder
	header, documentList string
}

// header struct that store terms descriptions and indices count
type header struct {
	Indices uint32
	Terms   []termDescription
}

// termDescription stores term, indice, postingList size and postingList file position
type termDescription struct {
	Term                Term
	Indice              uint32
	PostingListSize     uint32
	PostingListPosition uint32
}

// Load loads inverted index structure from disk
func (r *onDiscIndicesReader) Load() (InvertedIndexIndices, error) {
	// read header
	headerFile, err := os.Open(r.header)
	if err != nil {
		return nil, err
	}

	defer headerFile.Close()
	var header = new(header)

	decoder := gob.NewDecoder(headerFile)
	if err = decoder.Decode(header); err != nil {
		return nil, err
	}

	// mmap document list file
	documentListFile, err := os.Open(r.documentList)
	if err != nil {
		return nil, err
	}

	defer documentListFile.Close()
	data, err := mmap.Map(documentListFile, mmap.RDONLY, 0)

	if err != nil {
		return nil, err
	}

	// create inverted index structure slice
	indices := make([]InvertedIndex, int(header.Indices))
	invertedIndexStructureIndices := make([]invertedIndexStructure, len(indices))

	// here we create list of invertedIndexStructure
	for _, description := range header.Terms {
		if description.PostingListSize == 0 {
			invertedIndexStructureIndices[description.Indice] = nil
			continue
		}

		if invertedIndexStructureIndices[description.Indice] == nil {
			invertedIndexStructureIndices[description.Indice] = make(invertedIndexStructure)
		}

		invertedIndexStructureIndices[description.Indice][description.Term] = struct {
			size     uint32
			position uint32
		}{size: description.PostingListSize, position: description.PostingListPosition}
	}

	// create NewOnDiscInvertedIndex for given indice
	for i, invertedIndexStructure := range invertedIndexStructureIndices {
		if invertedIndexStructure == nil {
			indices[i] = nil
		} else {
			indices[i] = NewOnDiscInvertedIndex(data, r.decoder, invertedIndexStructure)
		}
	}

	invertedIndices := NewInvertedIndexIndices(indices)

	// we should unmap data when invertedIndices will be destroyed
	runtime.SetFinalizer(invertedIndices, func(r *invertedIndexIndicesImpl) {
		data.Unmap()
	})

	return invertedIndices, nil
}

// onDiscIndicesWriter implements IndicesWriter interface
type onDiscIndicesWriter struct {
	encoder              compression.Encoder
	header, documentList string
}

// Save tries to save index on disc, return non nil error on failure
func (w *onDiscIndicesWriter) Save(indices Indices) error {
	// open header list
	headerFile, err := os.OpenFile(
		w.header,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}

	defer headerFile.Close()

	// open document list
	documentListFile, err := os.OpenFile(
		w.documentList,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}

	defer documentListFile.Close()
	// use buffered writer for providing efficient writing to file
	documentListBuf := bufio.NewWriter(documentListFile)

	// mapValueOffset store current posting list offset
	mapValueOffset := int64(0)

	// header struct that should be loaded on Load
	header := header{
		Terms:   make([]termDescription, 0),
		Indices: uint32(len(indices)),
	}

	for indice, index := range indices {
		if index == nil {
			continue
		}

		for term, postingList := range index {
			// there is not possible, we should throw error
			if postingList == nil {
				return PostingListShouldBeNotNilError
			}

			// Encode posting list into value
			value := w.encoder.Encode(postingList)

			header.Terms = append(header.Terms, termDescription{
				Term:                term,
				Indice:              uint32(indice),
				PostingListSize:     uint32(len(value)),
				PostingListPosition: uint32(mapValueOffset),
			})

			if _, err = documentListBuf.Write(value); err != nil {
				return err
			}

			mapValueOffset += int64(len(value))
		}
	}

	encoder := gob.NewEncoder(headerFile)
	encoder.Encode(header)

	err = documentListBuf.Flush()
	if err != nil {
		return err
	}

	return nil
}
