package suggest

import (
	"bufio"
	"encoding/binary"
	"io"
)

// InvertedIndexIndicesReader represents entity for loading InvertedIndexIndices from storage
type InvertedIndexIndicesReader interface {
	// Load loads inverted index indices structure from disk
	Load() (InvertedIndexIndices, error)
}

// InvertedIndexIndicesWriter represents entity for saving Index structure in storage
type InvertedIndexIndicesWriter interface {
	// Save tries to save index on disc, return non nil error on failure
	Save(index Index) error
}

// NewOnDiscInvertedIndexWriter returns new instance of InvertedIndexIndicesWriter
func NewOnDiscInvertedIndexWriter(encoder Encoder, header io.WriteSeeker, documentList io.Writer, fromPosition int64) InvertedIndexIndicesWriter {
	return &onDiscWriter{
		encoder:      encoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

// NewOnDiscInvertedIndexReader returns new instance of InvertedIndexIndicesReader
func NewOnDiscInvertedIndexReader(decoder Decoder, header io.ReaderAt, documentList io.ReaderAt, fromPosition int64) InvertedIndexIndicesReader {
	return &onDiscReader{
		decoder:      decoder,
		header:       header,
		documentList: documentList,
		fromPosition: fromPosition,
	}
}

// onDiscReader implements InvertedIndexIndicesReader interface
type onDiscReader struct {
	decoder      Decoder
	header       io.ReaderAt
	documentList io.ReaderAt
	fromPosition int64
}

// Load loads inverted index indices structure from disk
func (r *onDiscReader) Load() (InvertedIndexIndices, error) {
	buf := make([]byte, 4)

	// first of all we read indices length
	_, err := r.header.ReadAt(buf, r.fromPosition)
	if err != nil {
		return nil, err
	}

	indicesLength := binary.LittleEndian.Uint32(buf)
	// create indices struct
	indices := make([]InvertedIndex, indicesLength)
	// each indices slot represents map length, map position each uint32
	buf = make([]byte, indicesLength*8)
	r.header.ReadAt(buf, r.fromPosition+4)

	for i := uint32(0); i < indicesLength; i++ {
		k := i * 8
		position, length := binary.LittleEndian.Uint32(buf[k:k+4]), binary.LittleEndian.Uint32(buf[k+4:k+8])

		// length == 0 marks that there is not inverted index for given length
		if length == 0 {
			indices[i] = nil
			continue
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

		indices[i] = NewOnDiscInvertedIndex(r.documentList, r.decoder, m)
	}

	return NewInvertedIndexIndices(indices), nil
}

// onDiscWriter implements InvertedIndexIndicesWriter interface
type onDiscWriter struct {
	encoder      Encoder
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
func (w *onDiscWriter) Save(index Index) error {
	// Seek header to fromPosition
	_, err := w.header.Seek(w.fromPosition, io.SeekStart)
	if err != nil {
		return err
	}

	headerBuf := bufio.NewWriter(w.header)
	documentListBuf := bufio.NewWriter(w.documentList)

	// Save indices length in header (4 bytes)
	err = binary.Write(headerBuf, binary.LittleEndian, uint32(len(index)))
	if err != nil {
		return err
	}

	// mapOffset
	mapOffset := w.fromPosition + 4 + int64(8*len(index))
	// mapValueOffset store posting lists
	mapValueOffset := int64(0)

	for _, table := range index {
		if table == nil {
			// if there is no table, store 0 as inverted index offset
			err = writePair(headerBuf, 0, 0)
			if err != nil {
				return err
			}

			continue
		}

		// otherwise store map structure offset + map size
		writePair(headerBuf, uint32(mapOffset), uint32(len(table)))
		// map size * (term + posting list bytes size + mapValueOffset)
		mapOffset += int64(len(table) * 12)
	}

	for _, table := range index {
		if table == nil {
			continue
		}

		for term, postingList := range table {
			// there is not possible, we should throw error
			if postingList == nil {
				panic("postingList is nil in table")
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
