package language_model

type WordCount = uint32

type TrieWalker = func(path []WordId, count WordCount)

// Trie represents data structure for counting ngrams.
type Trie interface {
	// Find returns WordCount for given wordId sequence.
	Find(sentence []WordId) (WordCount, error)
	// Put increments WordCount for last element of given sequence.
	Put(sentence []WordId) error
	// Walk iterates through trie and calls walker function on each element.
	Walk(walker TrieWalker)
}

// NewTrie creates new instance of Trie
func NewTrie() *trie {
	return &trie{
		root: &node{
			children: make(childrenTable),
			count:    0,
		},
	}
}

// trie implements Trie data structure
type trie struct {
	root *node
}

// node represents trie element
type node struct {
	children childrenTable
	count    WordCount
}

// childrenTable
type childrenTable map[WordId]*node

func (t *trie) Find(sentence []WordId) (WordCount, error) {
	n := t.root

	for _, w := range sentence {
		n = n.children[w]

		if n == nil {
			break
		}
	}

	if n != nil {
		return n.count, nil
	}

	return 0, nil
}

func (t *trie) Put(sentence []WordId) error {
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

func (t *trie) Walk(walker TrieWalker) {
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
