package lm

import (
	"errors"
	"log"
	"sort"
)

type (
	// ContextOffset represents the id of parent nGram path
	ContextOffset = uint32
	key           = uint64
)

// NGramVector represents one level of nGram trie
type NGramVector interface {
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
	keys   []key
	values []WordCount
	total  WordCount
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
func makeKey(word WordID, context ContextOffset) key {
	if context > maxContextOffset {
		log.Fatal(errors.New("Out of maxContextOffset"))
	}

	return pack(context, word)
}

// getWordID returns the word id for the given key
func getWordID(key key) WordID {
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
