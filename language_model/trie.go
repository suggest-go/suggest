package language_model

const maxOrder = 5

type Trie interface {
	Find(sentence []wordId) (uint32, error)
	Put(sentence []wordId) error
	Walk(observer func(path []wordId, count uint32))
}

func NewTrie() *trie {
	return &trie{
		root: &node{
			children: make(map[uint32]*node),
			count:    0,
		},
	}
}

type trie struct {
	root *node
}

type node struct {
	children map[uint32]*node
	count    uint32
}

func (t *trie) Find(sentence []wordId) (uint32, error) {
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

func (t *trie) Put(sentence []wordId) error {
	n := t.root

	for _, w := range sentence {
		child := n.children[w]

		if child == nil {
			if n.children == nil {
				n.children = make(map[uint32]*node)
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

func (t *trie) Walk(observer func(path []wordId, count uint32)) {
	t.root.walk([]wordId{}, observer)
}

func (n *node) walk(path []wordId, observer func(path []wordId, count uint32)) {
	observer(path, n.count)

	if n.children == nil {
		return
	}

	for w, child := range n.children {
		child.walk(append(path, w), observer)
	}
}
