package suggest

import (
	"io"
	"encoding/binary"
)

type InvertedIndexIndicesReader interface {
	Load() InvertedIndexIndices
}

type InvertedIndexIndicesWriter interface {
	Save(index Index) error
}

func NewOnDiscInvertedIndexWriter(encoder Encoder, header io.WriteSeeker, documentList io.Writer, fromPosition int64) InvertedIndexIndicesWriter {
	return &onDiscWriter{encoder, header, documentList, fromPosition}
}

type onDiscWriter struct {
	encoder Encoder
	header  io.WriteSeeker
	documentList io.Writer
	fromPosition int64
}

// *** HEADER ***
// INDICES_LENGTH - 4 byte
// [<pos1, length1>, ..., <posN, lengthN>] - 4 * 8 byte
// <pos1, length1>:
// pos1 - <term1, postingListPos>
// pos2 - <term2, postingListPos>
// ...
// posN - <termN, postingListPos>
//
// *** DOC LIST ***
// Values
func (w *onDiscWriter) Save(index Index) error {
	// Seek header to fromPosition
	_, err := w.header.Seek(w.fromPosition, io.SeekStart)
	if err != nil {
		return err
	}

	// Save indices length in header (4 bytes)
	err = binary.Write(w.header, binary.LittleEndian, uint32(len(index)))
	if err != nil {
		return err
	}

	mapOffset := w.fromPosition + 4 + int64(8 * len(index))
	// mapValueOffset store posting lists
	mapValueOffset := int64(0)

	for _, table := range index {
		if table == nil {
			// if there is no table, store 0 as inverted index offset
			err = writePair(w.header, 0, 0)
			if err != nil {
				return err
			}

			continue
		}

		// otherwise store map structure offset + map size
		writePair(w.header, uint32(mapOffset), uint32(len(table)))
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
			err = writeTro(w.header, uint32(term), uint32(len(value)), uint32(mapValueOffset))

			if err != nil {
				return err
			}

			_, err = w.documentList.Write(value)
			if err != nil {
				return err
			}

			mapValueOffset += int64(len(value))
		}
	}

	return nil
}

// writePair writes binary representation of two uint32 numbers to io.Writer
func writePair(writer io.Writer, a, b uint32) error {
	var pairBuf = []uint32{a, b}
	return binary.Write(writer, binary.LittleEndian, pairBuf)
}

// writeTro writes binary representation of 3 uint32 numbers to io.Writer
func writeTro(writer io.Writer, a, b, c uint32) error {
	var buf = []uint32{a, b, c}
	return binary.Write(writer, binary.LittleEndian, buf)
}

type onDiscReader struct {
	decoder Decoder
	header io.ReaderAt
	documentList io.ReaderAt
	fromPosition int64
}

func (r *onDiscReader) Load() InvertedIndexIndices {
	buf := make([]byte, 4)
	headerPos := r.fromPosition

	_, err := r.header.ReadAt(buf, headerPos)
	if err != nil {
		panic(err)
	}

	headerPos += 4

	indicesLength := binary.LittleEndian.Uint32(buf)
	indices := make([]InvertedIndex, indicesLength)
	buf = make([]byte, indicesLength * 8)
	r.header.ReadAt(buf, headerPos)
	headerPos += int64(8 * indicesLength)

	for i := uint32(0); i < indicesLength; i++ {
		k := i * 8
		position, length := binary.LittleEndian.Uint32(buf[k:k+4]), binary.LittleEndian.Uint32(buf[k+4:k+8])

		if length == 0 {
			indices[i] = nil
			continue
		}


		mapBuf := make([]byte, 12 * length)
		r.header.ReadAt(mapBuf, int64(position))

		m := make(invertedIndexStructure, length)
		var term, size, offset uint32

		for l := uint32(0); l < length; l++ {
			j := l * 12
			term = binary.LittleEndian.Uint32(mapBuf[j:j+4])
			size = binary.LittleEndian.Uint32(mapBuf[j+4:j+8])
			offset = binary.LittleEndian.Uint32(mapBuf[j+8:j+12])

			m[Term(term)] = struct {
				size     uint32
				position uint32
			}{size: size, position: offset}
		}

		indices[i] = &onDiscInvertedIndex{r.documentList, r.decoder, m}
	}

	return NewInvertedIndexIndices(indices)
}

type invertedIndexStructure map[Term]struct{ size uint32; position uint32 }

type onDiscInvertedIndex struct {
	reader io.ReaderAt
	decoder Decoder
	m invertedIndexStructure
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
