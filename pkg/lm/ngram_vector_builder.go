package lm

import (
	"errors"
	"fmt"

	"github.com/alldroll/go-datastructures/rbtree"
)

// NGramVectorBuilder is an entity that responses for building NGramVector
type NGramVectorBuilder interface {
	// Put adds the given sequence of nGrams and count to model
	Put(nGrams []WordID, count WordCount) error
	// Build creates new instance of NGramVector
	Build() NGramVector
}

// NGramVectorFactory represents a factory method for creating a NGramVector instance.
type NGramVectorFactory func(tree rbtree.Tree) NGramVector

// ErrNGramOrderIsOutOfRange informs that the given NGrams is out of range for the given
var ErrNGramOrderIsOutOfRange = errors.New("nGrams order is out of range")

// nGramVectorBuilder implements NGramVectorBuilder interface
type nGramVectorBuilder struct {
	parents []NGramVector
	factory NGramVectorFactory
	tree    rbtree.Tree
}

// NewNGramVectorBuilder creates new instance of NGramVectorBuilder
func NewNGramVectorBuilder(parents []NGramVector, factory NGramVectorFactory) NGramVectorBuilder {
	return &nGramVectorBuilder{
		parents: parents,
		factory: factory,
		tree:    rbtree.New(),
	}
}

// Put adds the given sequence of nGrams and count to model
func (m *nGramVectorBuilder) Put(nGrams []WordID, count WordCount) error {
	if len(nGrams) != len(m.parents)+1 {
		return ErrNGramOrderIsOutOfRange
	}

	parent := InvalidContextOffset

	for i, nGram := range nGrams {
		if i == len(nGrams)-1 {
			node := &nGramNode{
				key:   makeKey(nGram, parent),
				value: count,
			}

			prev := m.tree.Find(node)

			if prev != nil {
				(prev.(*nGramNode)).value += count
			} else {
				if _, err := m.tree.Insert(node); err != nil {
					return fmt.Errorf("failed to insert the node: %w", err)
				}
			}
		} else {
			parent = m.parents[i].GetContextOffset(nGram, parent)
		}
	}

	return nil
}

// Build creates new instance of NGramVector
func (m *nGramVectorBuilder) Build() NGramVector {
	return m.factory(m.tree)
}

// nGramNode represents tree node for the given nGram
type nGramNode struct {
	key   uint64
	value WordCount
}

// Less tells is current elements is bigger than the other
func (n *nGramNode) Less(other rbtree.Item) bool {
	return n.key < other.(*nGramNode).key
}
