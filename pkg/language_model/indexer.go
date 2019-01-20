package language_model

import (
	"errors"
	"sync"
)

type WordId = uint32

type Indexer interface {
	//
	GetOrCreate(nGram Word) WordId
	//
	Find(id WordId) (Word, error)
}

func NewIndexer() *indexer {
	return &indexer{
		table:  map[Word]WordId{},
		holder: []Word{},
		lock:   sync.RWMutex{},
	}
}

type indexer struct {
	table  map[Word]WordId
	holder []Word
	lock   sync.RWMutex
}

func (i *indexer) GetOrCreate(nGram Word) WordId {
	i.lock.RLock()
	index, ok := i.table[nGram]
	i.lock.RUnlock()

	if !ok {
		i.lock.Lock()
		defer i.lock.Unlock()

		index, ok = i.table[nGram]
		if ok {
			return index
		}

		index = uint32(len(i.holder))
		i.table[nGram] = index
		i.holder = append(i.holder, nGram)
	}

	return index
}

func (i *indexer) Find(index WordId) (Word, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	if uint32(len(i.holder)) <= index {
		return "", errors.New("index is not exists")
	}

	return i.holder[int(index)], nil
}
