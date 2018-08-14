package index

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/alldroll/suggest/compression"
	mmap "github.com/edsrzf/mmap-go"
	"io"
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

// Load loads inverted index structure from disk
func (r *onDiscIndicesReader) Load() (InvertedIndexIndices, error) {
	headerFile, err := os.Open(r.header)
	if err != nil {
		return nil, err
	}

	headerBuf := headerFile //bufio.NewReader(headerFile)
	defer headerFile.Close()

	buf := make([]byte, 4)

	// first of all we read indices length
	_, err = headerBuf.Read(buf)
	if err != nil {
		return nil, err
	}

	indicesLength := binary.LittleEndian.Uint32(buf)
	// create indices struct
	indices := make([]InvertedIndex, indicesLength)

	// read indices structure
	buf = make([]byte, 4*2*indicesLength)
	_, err = headerBuf.Read(buf)
	if err != nil {
		return nil, err
	}

	documentListFile, err := os.Open(r.documentList)
	if err != nil {
		return nil, err
	}

	defer documentListFile.Close()

	data, err := mmap.Map(documentListFile, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < indicesLength; i++ {
		k := i * 8
		position := binary.LittleEndian.Uint32(buf[k:])
		length := binary.LittleEndian.Uint32(buf[k+4:])

		if length == 0 {
			indices[i] = nil
			continue
		}

		var (
			term, size, offset, j uint32
			// m
			m = make(invertedIndexStructure, length)
			// each map entry represents key (term) and 2 uint32 (posting list byte size and position in documentList file)
			// so we have 3 uint32 numbers
			mapBuf = make([]byte, 12*length)
		)

		_, err = headerBuf.ReadAt(mapBuf, int64(position))
		if err != nil {
			return nil, err
		}

		for l := uint32(0); l < length; l++ {
			j = l * 12
			term = binary.LittleEndian.Uint32(mapBuf[j:])
			size = binary.LittleEndian.Uint32(mapBuf[j+4:])
			offset = binary.LittleEndian.Uint32(mapBuf[j+8:])

			m[Term(term)] = struct {
				size     uint32
				position uint32
			}{size: size, position: offset}
		}

		indices[i] = NewOnDiscInvertedIndex(data, r.decoder, m)
	}

	invertedIndices := NewInvertedIndexIndices(indices)

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
//
// *** HEADER FILE ***
// INDICES_LENGTH - 4 byte
// [<pos1, length1>, ..., <posN, lengthN>] - 4 * 8 byte
//
// <pos1, length1>:
// <term11, size, postingListPos>
// <term12, size, postingListPos>
// <term13, size, postingListPos>
// <term1length1, size, postingListPos>
// ...
// posN - <termN, size, postingListPos>
// <termN1, size, postingListPos>
// <termN2, size, postingListPos>
// <termN3, size, postingListPos>
// <termNlengthN, size, postingListPos>

// *** DOCLIST FILE***
// Binary encoded values
//
func (w *onDiscIndicesWriter) Save(indices Indices) error {
	headerFile, err := os.OpenFile(
		w.header,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return err
	}

	headerBuf := bufio.NewWriter(headerFile)
	defer headerFile.Close()

	documentListFile, err := os.OpenFile(
		w.documentList,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return err
	}

	documentListBuf := bufio.NewWriter(documentListFile)
	defer documentListFile.Close()

	// Save indices length in header (4 bytes)
	err = binary.Write(headerBuf, binary.LittleEndian, uint32(len(indices)))
	if err != nil {
		return err
	}

	// mapOffset store inverted index structure
	mapOffset := 4 + int64(8*len(indices))
	// mapValueOffset store posting lists
	mapValueOffset := int64(0)

	for _, index := range indices {
		if index == nil {
			// if there is no table, store 0 as inverted index offset
			err = writePair(headerBuf, 0, 0)
			if err != nil {
				return err
			}

			continue
		}

		// otherwise store map structure offset + map size
		err = writePair(headerBuf, uint32(mapOffset), uint32(len(index)))
		if err != nil {
			return err
		}

		// map size * (term + posting list bytes size + mapValueOffset)
		mapOffset += int64(len(index) * 12)
	}

	for _, index := range indices {
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
			// Write <term, posting list size, posting list offset)
			err = writeTrio(headerBuf, uint32(term), uint32(len(value)), uint32(mapValueOffset))

			if err != nil {
				return err
			}

			_, err = documentListBuf.Write(value)
			if err != nil {
				return err
			}

			mapValueOffset += int64(len(value))
		}
	}

	err = headerBuf.Flush()
	if err != nil {
		return err
	}

	err = documentListBuf.Flush()
	if err != nil {
		return err
	}

	return nil
}

// writePair writes binary representation of two uint32 numbers to io.Writer
func writePair(writer io.Writer, a, b uint32) error {
	var buf = []uint32{a, b}
	return binary.Write(writer, binary.LittleEndian, buf)
}

// writeTrio writes binary representation of 3 uint32 numbers to io.Writer
func writeTrio(writer io.Writer, a, b, c uint32) error {
	var buf = []uint32{a, b, c}
	return binary.Write(writer, binary.LittleEndian, buf)
}
