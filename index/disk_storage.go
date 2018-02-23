package index

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/alldroll/suggest/compression"
	"io"
)

var (
	PostingListShouldBeNotNilError = errors.New("PostingList should be not nil")
)

// NewOnDiscIndicesReader returns new instance of IndexReader
func NewOnDiscIndexReader(decoder compression.Decoder, header io.ReaderAt, documentList io.ReaderAt, fromPosition int64) IndexReader {
	return &onDiscIndexReader{
		decoder:      decoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

// NewOnDiscIndicesReader returns new instance of IndicesReader
func NewOnDiscIndicesReader(decoder compression.Decoder, header io.ReaderAt, documentList io.ReaderAt, fromPosition int64) IndicesReader {
	return &onDiscIndicesReader{
		decoder:      decoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

// NewOnDiscIndexWriter returns new instance of IndexWriter
func NewOnDiscIndexWriter(encoder compression.Encoder, header io.WriteSeeker, documentList io.Writer, fromPosition int64) IndexWriter {
	return &onDiscIndexWriter{
		encoder:      encoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

// NewOnDiscIndicesWriter returns new instance of IndicesWriter
func NewOnDiscIndicesWriter(encoder compression.Encoder, header io.WriteSeeker, documentList io.Writer, fromPosition int64) IndicesWriter {
	return &onDiscIndicesWriter{
		encoder:      encoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

type onDiscIndexReader struct {
	decoder      compression.Decoder
	header       io.ReaderAt
	documentList io.ReaderAt
	fromPosition int64
}

// Load loads inverted index structure from disk
func (r *onDiscIndexReader) Load() (InvertedIndex, error) {
	buf := make([]byte, 8)
	// first of all we read length and position of inverted index
	_, err := r.header.ReadAt(buf, r.fromPosition)
	if err != nil {
		return nil, err
	}

	position := binary.LittleEndian.Uint32(buf[:4])
	length := binary.LittleEndian.Uint32(buf[4:])

	if length == 0 {
		return nil, nil
	}

	var (
		term, size, offset, j uint32
		m                     = make(invertedIndexStructure, length)
		// each map entry represents key (term) and 2 uint32 (posting list byte size and position in documentList file)
		// so we have 3 uint32 numbers
		mapBuf = make([]byte, 12*length)
	)

	_, err = r.header.ReadAt(mapBuf, int64(position))
	if err != nil {
		return nil, err
	}

	for l := uint32(0); l < length; l++ {
		j = l * 12
		term = binary.LittleEndian.Uint32(mapBuf[j : j+4])
		size = binary.LittleEndian.Uint32(mapBuf[j+4 : j+8])
		offset = binary.LittleEndian.Uint32(mapBuf[j+8 : j+12])

		m[Term(term)] = struct {
			size     uint32
			position uint32
		}{size: size, position: offset}
	}

	return NewOnDiscInvertedIndex(r.documentList, r.decoder, m), nil
}

// onDiscIndicesReader implements Reader interface
type onDiscIndicesReader struct {
	decoder      compression.Decoder
	header       io.ReaderAt
	documentList io.ReaderAt
	fromPosition int64
}

// Load loads inverted index structure from disk
func (r *onDiscIndicesReader) Load() (InvertedIndexIndices, error) {
	buf := make([]byte, 4)

	// first of all we read indices length
	_, err := r.header.ReadAt(buf, r.fromPosition)
	if err != nil {
		return nil, err
	}

	position := r.fromPosition + 4
	indicesLength := binary.LittleEndian.Uint32(buf)
	// create indices struct
	indices := make([]InvertedIndex, indicesLength)

	for i := uint32(0); i < indicesLength; i++ {
		indexReader := NewOnDiscIndexReader(r.decoder, r.header, r.documentList, position)
		position += 8

		invertedIndex, err := indexReader.Load()
		if err != nil {
			return nil, err
		}

		indices[i] = invertedIndex
	}

	return NewInvertedIndexIndices(indices), nil
}

// onDiscIndexWriter implements IndexWriter interface
type onDiscIndexWriter struct {
	encoder      compression.Encoder
	header       io.WriteSeeker
	documentList io.Writer
	fromPosition int64
}

// Save tries to save index on disc, return non nil error on failure
//
// *** HEADER ***
// pos, length - 8 byte
// <term1, size, postingListPos>
// <term2, size, postingListPos>
// ...
// <termN, size, postingListPos>
//
// *** DOC LIST ***
// Values
//
func (w *onDiscIndexWriter) Save(index Index) error {
	// Seek header to fromPosition
	_, err := w.header.Seek(w.fromPosition, io.SeekStart)
	if err != nil {
		return err
	}

	headerBuf := bufio.NewWriter(w.header)
	documentListBuf := bufio.NewWriter(w.documentList)

	// mapOffset
	mapOffset := w.fromPosition + 8
	// mapValueOffset store posting lists
	mapValueOffset := int64(0)

	// store map structure offset + map size
	err = writePair(headerBuf, uint32(mapOffset), uint32(len(index)))
	if err != nil {
		return err
	}

	// map size * (term + posting list bytes size + mapValueOffset)
	mapOffset += int64(len(index) * 12)

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

// onDiscIndicesWriter implements IndicesWriter interface
type onDiscIndicesWriter struct {
	encoder      compression.Encoder
	header       io.WriteSeeker
	documentList io.Writer
	fromPosition int64
}

// Save tries to save index on disc, return non nil error on failure
//
// *** HEADER ***
// INDICES_LENGTH - 4 byte
// [<pos1, length1>, ..., <posN, lengthN>] - 4 * 8 byte
// <pos1, length1>:
// pos1 - <term1, size, postingListPos>
// pos2 - <term2, size, postingListPos>
// ...
// posN - <termN, size, postingListPos>
//
// *** DOC LIST ***
// Values
//
func (w *onDiscIndicesWriter) Save(indices Indices) error {
	// Seek header to fromPosition
	_, err := w.header.Seek(w.fromPosition, io.SeekStart)
	if err != nil {
		return err
	}

	headerBuf := bufio.NewWriter(w.header)
	documentListBuf := bufio.NewWriter(w.documentList)

	// Save indices length in header (4 bytes)
	err = binary.Write(headerBuf, binary.LittleEndian, uint32(len(indices)))
	if err != nil {
		return err
	}

	// mapOffset
	mapOffset := w.fromPosition + 4 + int64(8*len(indices))
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
