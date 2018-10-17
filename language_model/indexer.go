package language_model

import "errors"

type wordId = uint32

type Indexer interface {
	GetOrCreate(nGram Word) wordId
	Find(id wordId) (Word, error)
}

func NewIndexer() *indexer {
	return &indexer{
		table:  map[Word]wordId{},
		holder: []Word{},
	}
}

type indexer struct {
	table  map[Word]wordId
	holder []Word
}

func (i *indexer) GetOrCreate(nGram Word) wordId {
	index, ok := i.table[nGram]

	if !ok {
		index = uint32(len(i.holder))
		i.table[nGram] = index
		i.holder = append(i.holder, nGram)
	}

	return index
}

func (i *indexer) Find(index wordId) (Word, error) {
	if uint32(len(i.holder)) <= index {
		return "", errors.New("index is not exists")
	}

	return i.holder[int(index)], nil
}
