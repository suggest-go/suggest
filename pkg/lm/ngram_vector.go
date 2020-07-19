package lm

import (
	"errors"
	"math"
)

type (
	// ContextOffset represents the id of parent nGram path
	ContextOffset = uint32
	key           = uint64
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
