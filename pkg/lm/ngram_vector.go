package lm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/suggest-go/suggest/pkg/utils"
)

type (
	// ContextOffset represents the id of parent nGram path
	ContextOffset = uint32
	key = uint64
)

// NGramVector represents one level of nGram trie
type NGramVector interface {
	// GetCount returns WordCount and Node ContextOffset for the given pair (word, context)
	GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset)
	// GetContextOffset returns the given node context offset
	GetContextOffset(word WordID, context ContextOffset) ContextOffset
	// CorpusCount returns size of all counts in the collection
	CorpusCount() WordCount
	// SubVector returns NGramVector for the given context
	SubVector(context ContextOffset) NGramVector
}

const (
	// InvalidContextOffset is context id that represents invalid context offset
	InvalidContextOffset = maxContextOffset - 1
	maxUint32            = uint32(math.MaxUint32)
	maxContextOffset     = maxUint32 - 1
)

var (
	// ErrContextOverflow tells that it was an attempt
	ErrContextOverflow = errors.New("out of maxContextOffset")
)

type sortedArray struct {
	keys   []key
	values []WordCount
	total  WordCount
}

// GetCount returns WordCount and Node ContextOffset for the given pair (word, context)
func (s *sortedArray) GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset) {
	key := makeKey(word, context)
	i := s.find(key)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return s.values[int(i)], i
}

// GetContextOffset returns the given node context offset
func (s *sortedArray) GetContextOffset(word WordID, context ContextOffset) ContextOffset {
	key := makeKey(word, context)

	return s.find(key)
}

// CorpusCount returns size of all counts in the collection
func (s *sortedArray) CorpusCount() WordCount {
	return s.total
}

// SubVector returns NGramVector for the given context
func (s *sortedArray) SubVector(context ContextOffset) NGramVector {
	minChild := makeKey(0, context)
	maxChild := makeKey(maxContextOffset-2, context)

	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= minChild })

	if i < 0 || i >= len(s.keys) {
		return nil
	}

	j := sort.Search(len(s.keys)-i, func(j int) bool { return s.keys[j+i] >= maxChild })

	return &sortedArray{
		keys:   s.keys[i:i+j],
		values: s.values[i:i+j],
		total:  s.total,
	}
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (s *sortedArray) MarshalBinary() ([]byte, error) {
	var result bytes.Buffer

	encodedKeys := make([]byte, len(s.keys)*binary.MaxVarintLen64)
	prevKey := uint64(0)
	keyEndPos := 0

	// performs delta encoding
	for _, el := range s.keys {
		keyEndPos += binary.PutUvarint(encodedKeys[keyEndPos:], el-prevKey)
		prevKey = el
	}

	valEndPos := len(s.values) * 4
	encodedValues := make([]byte, valEndPos)

	for i, el := range s.values {
		binary.LittleEndian.PutUint32(encodedValues[i*4:(i+1)*4], el)
	}

	// allocate buffer capacity
	result.Grow(keyEndPos + valEndPos + 4 + strconv.IntSize*2)

	// write header
	if _, err := fmt.Fprintln(&result, keyEndPos, valEndPos, s.total); err != nil {
		return nil, err
	}

	// write data
	result.Write(encodedKeys[:keyEndPos])
	result.Write(encodedValues)

	return result.Bytes(), nil
}

// UnmarshalBinary decodes the binary form
func (s *sortedArray) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	keySize, valSize := 0, 0

	if _, err := fmt.Fscanln(buf, &keySize, &valSize, &s.total); err != nil {
		return err
	}

	n := 0
	keyEndPos := 0
	encodedKeys := buf.Next(keySize)
	s.keys = make([]key, valSize/4)
	prev := uint64(0)

	// 0, 1, 3, 6, 7 -> 0, 1, 4, 10, 17 (delta decoding)
	for i := 0; i < len(s.keys); i++ {
		s.keys[i], n = binary.Uvarint(encodedKeys[keyEndPos:])
		s.keys[i] += prev
		prev = s.keys[i]
		keyEndPos += n
	}

	encodedValues := buf.Next(valSize)
	s.values = make([]WordCount, valSize/4)

	for i := 0; i < len(s.values); i++ {
		s.values[i] = binary.LittleEndian.Uint32(encodedValues[i*4 : (i+1)*4])
	}

	return nil
}

// find finds the given key in the collection. Returns ContextOffset if the key exists, otherwise returns InvalidContextOffset
func (s *sortedArray) find(key uint64) ContextOffset {
	if len(s.keys) == 0 || s.keys[0] > key || s.keys[len(s.keys) - 1] < key {
		return InvalidContextOffset
	}

	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= key })

	if i < 0 || i >= len(s.keys) || s.keys[i] != key {
		return InvalidContextOffset
	}

	return ContextOffset(i)
}

// makeKey creates uint64 key for the given pair (word, context)
func makeKey(word WordID, context ContextOffset) key {
	if context > maxContextOffset {
		log.Fatal(ErrContextOverflow)
	}

	return utils.Pack(context, word)
}

// getWordID returns the word id for the given key
func getWordID(key key) WordID {
	return utils.UnpackRight(key)
}

func init() {
	gob.Register(&sortedArray{})
}
