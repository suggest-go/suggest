package language_model

import (
	"errors"
	"log"
	"sort"
)

type ContextOffset = uint32

//
type NGramVector interface {
	//
	Put(word, context WordId, count WordCount)

	//
	GetCount(word WordId, context ContextOffset) (WordCount, ContextOffset)

	//
	GetContextOffset(word WordId, context ContextOffset)

	//
	CorpousCount() WordCount
}

const (
	maxContextOffset     = uint32(0xffffffff)
	InvalidContextOffset = maxContextOffset - 1
)

type sortedArray struct {
	keys   []uint64
	values []uint32
	total  uint32
}

func NewNGramVector() *sortedArray {
	return &sortedArray{
		keys:   make([]uint64, 0),
		values: make([]uint32, 0),
		total:  0,
	}
}

//
func (s *sortedArray) Put(word WordId, context ContextOffset, count WordCount) {
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

//
func (s *sortedArray) GetCount(word WordId, context ContextOffset) (WordCount, ContextOffset) {
	key := makeKey(word, context)
	i := s.find(key)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return s.values[int(i)], i
}

//
func (s *sortedArray) CorpousCount() WordCount {
	return s.total
}

//
func (s *sortedArray) GetContextOffset(word WordId, context ContextOffset) ContextOffset {
	key := makeKey(word, context)
	return s.find(key)
}

//
func (s *sortedArray) find(key uint64) ContextOffset {
	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= key })

	if i < 0 || i >= len(s.keys) || s.keys[i] != key {
		return InvalidContextOffset
	}

	return ContextOffset(i)
}

func makeKey(word WordId, context ContextOffset) uint64 {
	if context > maxContextOffset {
		log.Fatal(errors.New("Out of maxContextOffset"))
	}

	return pack(uint32(context), word)
}

// pack packs 2 uint32 in uint64
func pack(a, b uint32) uint64 {
	return (uint64(a) << 32) | uint64(b)
}

// unpack explode uint64 into 2 uint32
func unpack(v uint64) (uint32, uint32) {
	return uint32(v >> 32), uint32(v & 0xFFFFFF)
}
