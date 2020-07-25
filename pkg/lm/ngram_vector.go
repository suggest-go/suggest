package lm

import (
	"errors"
	"log"
	"math"

	"github.com/suggest-go/suggest/pkg/store"
	"github.com/suggest-go/suggest/pkg/utils"
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
	// Store saves the given NGramVector to the provided output.
	Store(out store.Output) (int, error)
	// Load loads the saved NGramVector from the given in.
	Load(in store.Input) (int, error)
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

// getContext returns the context for the given key
func getContext(key key) WordID {
	return utils.UnpackLeft(key)
}
