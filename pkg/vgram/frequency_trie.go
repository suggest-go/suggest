package vgram

import "sort"

type FrequencyTrie interface {
	Find(gram string) Node
	Add(gram string)
	Walk(walker func(key string, node Node))
	Prune(threshold uint32)
}

type Node interface {
	frequencyHolder
	GetMarker() Marker
}

type Marker interface {
	frequencyHolder
}

type frequencyHolder interface {
	GetFrequency() uint32
}

type trie struct {
	root *node
	qMin uint32
}

func NewFrequencyTrie(qMin uint32) FrequencyTrie {
	return &trie{
		root: newNode(),
		qMin: qMin,
	}
}

func (t *trie) Find(gram string) Node {
	cur := t.root

	for _, char := range gram {
		cur = cur.children[char]

		if cur == nil {
			break
		}
	}

	return cur
}

func (t *trie) Add(gram string) {
	cur := t.root

	for i, char := range gram {
		child := cur.children[char]

		if child == nil {
			child = newNode()
			cur.children[char] = child
		}

		child.frequency++
		cur = child

		if uint32(i+1) >= t.qMin && cur.marker == nil {
			cur.marker = newMarker()
		}
	}

	if cur.marker != nil {
		cur.marker.frequency++
	}
}

func (t *trie) Walk(walker func(key string, n Node)) {
	t.root.walk("", walker)
}

func (t *trie) Prune(threshold uint32) {
	t.root.prune(threshold)
}

type nodeList []*node

// Len is the number of elements in the collection.
func (p nodeList) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p nodeList) Less(i, j int) bool { return p[i].frequency < p[j].frequency }

// Swap swaps the elements with indexes i and j.
func (p nodeList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type node struct {
	children  map[rune]*node
	marker    *marker
	frequency uint32
}

type marker struct {
	frequency uint32
}

func (m *marker) GetFrequency() uint32 {
	return m.frequency
}

func newNode() *node {
	return &node{
		children:  make(map[rune]*node),
		marker:    nil,
		frequency: 0,
	}
}

func newMarker() *marker {
	return &marker{
		frequency: 0,
	}
}

func (n *node) walk(key string, walker func(key string, n Node)) {
	for char, child := range n.children {
		walker(key+string(char), child)
		child.walk(key+string(char), walker)
	}
}

func (n *node) GetFrequency() uint32 {
	return n.frequency
}

func (n *node) GetMarker() Marker {
	if n.marker == nil {
		return nil
	}

	return n.marker
}

func (n *node) removeChildren() {
	if len(n.children) > 0 {
		n.children = make(map[rune]*node)
	}
}

func (n *node) getChildren() []*node {
	children := make([]*node, 0, len(n.children))

	for _, child := range n.children {
		children = append(children, child)
	}

	return children
}

func (n *node) removeChild(child *node) bool {
	success := false

	for char, candidate := range n.children {
		if candidate == child {
			delete(n.children, char)
			success = true
			break
		}
	}

	return success
}

func (n *node) prune(threshold uint32) {
	if n.marker == nil {
		for _, child := range n.getChildren() {
			child.prune(threshold)
		}

		return
	}

	leaf := n.marker
	freq := n.frequency

	if freq <= threshold {
		n.removeChildren()
		leaf.frequency = freq
	} else {
		// Select a maximal subset of children of n (excluding L),
		// so that the summation of their freq values and
		// L.freq is still not greater than T;
		//
		// Add the freq values of these children to
		// that of L, and remove these children from n;
		//
		// FOR (each remaining child c of n excluding L)
		// CALL Prune(c, T); // recursive call

		children := nodeList(n.getChildren())
		sort.Sort(children)

		leafFreq := leaf.frequency
		for _, child := range children {
			if leafFreq+child.frequency <= threshold {
				leafFreq += child.frequency
				n.removeChild(child)
			} else {
				child.prune(threshold)
			}
		}

		leaf.frequency = leafFreq
	}
}
