package lm

// WordCount is a count of a corresponding path
type WordCount = uint32

// TrieIterator is callback that
type TrieIterator = func(path []WordID, count WordCount)

// CountTrie represents data structure for counting ngrams.
type CountTrie interface {
	// Put increments WordCount for last element of given sequence.
	Put(sentence []WordID, count WordCount) error
	// Walk iterates through trie and calls walker function on each element.
	Walk(walker TrieIterator)
}

// NewCountTrie creates new instance of CountTrie
func NewCountTrie() CountTrie {
	return &countTrie{
		root: &node{
			children: make(childrenTable),
			count:    0,
		},
		depth: 0,
	}
}

// countTrie implements Trie data structure
type countTrie struct {
	root  *node
	depth int
}

// node represents trie element
type node struct {
	children childrenTable
	count    WordCount
}

// childrenTable represents map for children of given node
type childrenTable map[WordID]*node

// Put increments WordCount for last element of given sequence.
func (t *countTrie) Put(sentence []WordID, count WordCount) error {
	if len(sentence) > t.depth {
		t.depth = len(sentence)
	}

	n := t.root

	for _, w := range sentence {
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
	return nil
}

// Walk iterates through trie and calls walker function on each element.
func (t *countTrie) Walk(walker TrieIterator) {
	if t.depth == 0 {
		return
	}

	path := make([]WordID, t.depth)
	t.root.iterate(0, path, walker)
}

// iterates through given depth and calls iterator on each path
func (n *node) iterate(depth int, path []WordID, iterator TrieIterator) {
	if n.count > 0 {
		iterator(path[:depth], n.count)
	}

	if n.children == nil {
		return
	}

	for w, child := range n.children {
		path[depth] = w
		child.iterate(depth+1, path, iterator)
	}
}
