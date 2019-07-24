package lm

import "errors"

// WordCount is a count of a corresponding path
type WordCount = uint32

// TrieIterator is a callback that is called for each path of the given trie
type TrieIterator = func(path Sentence, count WordCount) error

// CountTrie represents a data structure for counting ngrams.
type CountTrie interface {
	// Put increments WordCount for last element of given sequence.
	Put(sentence Sentence, count WordCount)
	// Walk iterates through trie and calls walker function on each element.
	Walk(walker TrieIterator) error
}

// ErrInvalidIndex tells that there is not data for the provided index index
var ErrInvalidIndex = errors.New("index is not exists")

// NewCountTrie creates new a instance of CountTrie
func NewCountTrie() CountTrie {
	return &countTrie{
		root: &node{
			children: make(childrenTable),
			count:    0,
		},
		depth:  0,
		table:  map[Token]uint32{},
		holder: []Token{},
	}
}

// countTrie implements a Trie data structure
type countTrie struct {
	root   *node
	depth  int
	table  map[Token]uint32
	holder []Token
}

// node represents a trie element
type node struct {
	children childrenTable
	count    WordCount
}

// childrenTable represents a map for children of the given node
type childrenTable map[uint32]*node

// Put increments WordCount for the last element of the given sequence.
func (t *countTrie) Put(sentence Sentence, count WordCount) {
	if len(sentence) > t.depth {
		t.depth = len(sentence)
	}

	n := t.root

	for _, word := range sentence {
		w := t.mapToUint32(word)
		child := n.children[w]

		if child == nil {
			if n.children == nil {
				n.children = make(childrenTable)
			}

			child = &node{
				children: nil,
				count:    0,
			}

			n.children[w] = child
		}

		n = child
	}

	n.count += count
}

// Walk iterates through the trie and calls the walker function on each element.
func (t *countTrie) Walk(walker TrieIterator) (err error) {
	if t.depth == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	path := make([]Token, t.depth)
	t.root.iterate(t, 0, path, walker)

	return err
}

// mapToUint32 maps the given token to a index value
func (t *countTrie) mapToUint32(token Token) uint32 {
	index, ok := t.table[token]

	if !ok {
		index = uint32(len(t.holder))
		t.table[token] = index
		t.holder = append(t.holder, token)
	}

	return index
}

// mapFromUint32 restores a token from the given index
func (t *countTrie) mapFromUint32(index uint32) (Token, error) {
	if uint32(len(t.holder)) <= index {
		return UnknownWordSymbol, ErrInvalidIndex
	}

	return t.holder[int(index)], nil
}

// iterate iterates through the given depth and calls the iterator on each path
func (n *node) iterate(trie *countTrie, depth int, path []Token, iterator TrieIterator) {
	if n.count > 0 {
		if err := iterator(path[:depth], n.count); err != nil {
			panic(err)
		}
	}

	if n.children == nil {
		return
	}

	for w, child := range n.children {
		token, err := trie.mapFromUint32(w)

		if err != nil {
			panic(err)
		}

		path[depth] = token
		child.iterate(trie, depth+1, path, iterator)
	}
}
