package lm

import (
	"errors"
	"log"
	"sort"
)

// ContextOffset represents the id of parent nGram path
type ContextOffset = uint32

// NGramVector represents one level of nGram trie
type NGramVector interface {
	// Puts the given pair(word, context) in the collection
	Put(word, context WordID, count WordCount)
	// Returns WordCount and Node ContextOffset for the given pair (word, context)
	GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset)
	// Returns the given node context offset
	GetContextOffset(word WordID, context ContextOffset) ContextOffset
	// Returns size of all counts in the collection
	CorpousCount() WordCount
	// Next returns next words for the given context
	Next(context ContextOffset) []WordID
}

const (
	maxUint32        = uint32(0xffffffff)
	maxContextOffset = maxUint32 - 1
	// InvalidContextOffset is context id that represents invalid context offset
	InvalidContextOffset = maxContextOffset - 1
)

type sortedArray struct {
	keys   []uint64
	values []uint32
	total  uint32
}

// NewNGramVector creates new instance of NGramVector
func NewNGramVector() NGramVector {
	return &sortedArray{
		keys:   make([]uint64, 0),
		values: make([]uint32, 0),
		total:  0,
	}
}

// Puts the given pair(word, context) in the collection
func (s *sortedArray) Put(word WordID, context ContextOffset, count WordCount) {
	key := makeKey(word, context)
	s.total += count

	switch i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= key }); {
	case i < 0:
		s.keys = append([]uint64{key}, s.keys...)
		s.values = append([]uint32{count}, s.values...)
	case i >= len(s.keys):
		s.keys = append(s.keys, key)
		s.values = append(s.values, count)
	default:
		if s.keys[i] == key {
			s.values[i] += count
		} else {
			s.keys = append(s.keys[:i], append([]uint64{key}, s.keys[i:]...)...)
			s.values = append(s.values[:i], append([]uint32{count}, s.values[i:]...)...)
		}
	}
}

// Returns WordCount and Node ContextOffset for the given pair (word, context)
func (s *sortedArray) GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset) {
	key := makeKey(word, context)
	i := s.find(key)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return s.values[int(i)], i
}

// Returns the given node context offset
func (s *sortedArray) GetContextOffset(word WordID, context ContextOffset) ContextOffset {
	key := makeKey(word, context)
	return s.find(key)
}

// Returns size of all counts in the collection
func (s *sortedArray) CorpousCount() WordCount {
	return s.total
}

// Next returns next words for the given context
func (s *sortedArray) Next(context ContextOffset) []WordID {
	minChild := makeKey(0, context)
	maxChild := makeKey(maxContextOffset-2, context)

	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= minChild })
	words := []WordID{}

	if i < 0 || i >= len(s.keys) {
		return words
	}

	for ; s.keys[i] <= maxChild; i++ {
		words = append(words, getWordID(s.keys[i]))
	}

	return words
}

// Finds given key in the collection. Returns ContextOffset if the given key exists
func (s *sortedArray) find(key uint64) ContextOffset {
	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= key })

	if i < 0 || i >= len(s.keys) || s.keys[i] != key {
		return InvalidContextOffset
	}

	return ContextOffset(i)
}

// Creates uint64 key for the given pair (word, context)
func makeKey(word WordID, context ContextOffset) uint64 {
	if context > maxContextOffset {
		log.Fatal(errors.New("Out of maxContextOffset"))
	}

	return pack(context, word)
}

// getWordID returns the word id for the given key
func getWordID(key uint64) uint32 {
	_, wordID := unpack(key)
	return wordID
}

// Packs 2 uint32 in uint64
func pack(a, b uint32) uint64 {
	return (uint64(a) << 32) | uint64(b&maxUint32)
}

// Unpacks explode uint64 into 2 uint32
func unpack(v uint64) (uint32, uint32) {
	return uint32(v >> 32), uint32(v)
}
