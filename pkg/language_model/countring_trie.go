package language_model

type WordCount = uint32

type TrieWalker = func(path []WordId, count WordCount)

// CountingTrie represents data structure for counting ngrams.
type CountingTrie interface {
	// Put increments WordCount for last element of given sequence.
	Put(sentence []WordId) error
	// Walk iterates through trie and calls walker function on each element.
	Walk(walker TrieWalker)
}

// NewCountingTrie creates new instance of CountingTrie
func NewTrie() *countingTrie {
	return &countingTrie{
		root: &node{
			children: make(childrenTable),
			count:    0,
		},
	}
}

// countingTrie implements Trie data structure
type countingTrie struct {
	root *node
}

// node represents trie element
type node struct {
	children childrenTable
	count    WordCount
}

// childrenTable represents map for children of given node
type childrenTable map[WordId]*node

func (t *countingTrie) Put(sentence []WordId) error {
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

	n.count++
	return nil
}

func (t *countingTrie) Walk(walker TrieWalker) {
	t.root.walk([]WordId{}, walker)
}

func (n *node) walk(path []WordId, walker TrieWalker) {
	walker(path, n.count)

	if n.children == nil {
		return
	}

	for w, child := range n.children {
		child.walk(append(path, w), walker)
	}
}
