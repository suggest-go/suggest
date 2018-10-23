package language_model

import "errors"

type WordId = uint32

type Indexer interface {
	GetOrCreate(nGram Word) WordId
	Find(id WordId) (Word, error)
}

func NewIndexer() *indexer {
	return &indexer{
		table:  map[Word]WordId{},
		holder: []Word{},
	}
}

type indexer struct {
	table  map[Word]WordId
	holder []Word
}

func (i *indexer) GetOrCreate(nGram Word) WordId {
	index, ok := i.table[nGram]

	if !ok {
		index = uint32(len(i.holder))
		i.table[nGram] = index
		i.holder = append(i.holder, nGram)
	}

	return index
}

func (i *indexer) Find(index WordId) (Word, error) {
	if uint32(len(i.holder)) <= index {
		return "", errors.New("index is not exists")
	}

	return i.holder[int(index)], nil
}
