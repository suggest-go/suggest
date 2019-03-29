package lm

import (
	"errors"
	"sync"
)

// WordID is an index of the corresponding word
type WordID = uint32

const (
	// UnknownWordID is an index of unregistered word
	UnknownWordID = uint32(0xffffffff)
)

// Indexer enumerates words in the vocabulary of a language model. Stores a two-way
// mapping between uint32 and strings.
type Indexer interface {
	// Gets the index for the word, creates if necessary.
	GetOrCreate(token Token) WordID
	// Returns the index for the word, otherwise returns UnknownWordID
	Get(token Token) WordID
	// Find a token by the given index
	Find(id WordID) (Token, error)
}

// NewIndexer creates new instance of Indexer
func NewIndexer() Indexer {
	return &indexer{
		table:  map[Token]WordID{},
		holder: []Token{},
		lock:   sync.RWMutex{},
	}
}

// indexer implements Indexer interface
type indexer struct {
	table  map[Token]WordID
	holder []Token
	lock   sync.RWMutex
}

// Gets the index for the word, creates if necessary.
func (i *indexer) GetOrCreate(token Token) WordID {
	i.lock.RLock()
	index, ok := i.table[token]
	i.lock.RUnlock()

	if !ok {
		i.lock.Lock()
		defer i.lock.Unlock()

		index, ok = i.table[token]
		if ok {
			return index
		}

		index = uint32(len(i.holder))
		i.table[token] = index
		i.holder = append(i.holder, token)
	}

	return index
}

// Returns the index for the word, otherwise returns UnknownWordID
func (i *indexer) Get(token Token) WordID {
	i.lock.RLock()
	index, ok := i.table[token]
	i.lock.RUnlock()

	if !ok {
		index = UnknownWordID
	}

	return index
}

// Find a token by the given index
// TODO replace error with <UNK> symbol
func (i *indexer) Find(index WordID) (Token, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	if uint32(len(i.holder)) <= index {
		return "", errors.New("index is not exists")
	}

	return i.holder[int(index)], nil
}
