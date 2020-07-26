package lm

import (
	"errors"
	"fmt"
	"log"

	"github.com/alldroll/go-datastructures/rbtree"
	"github.com/suggest-go/suggest/pkg/utils"
)

// NGramVectorBuilder is an entity that responses for building NGramVector
type NGramVectorBuilder interface {
	// Put adds the given sequence of nGrams and count to model
	Put(nGrams []WordID, count WordCount) error
	// Build creates new instance of NGramVector
	Build() NGramVector
}

// NGramVectorFactory represents a factory method for creating a NGramVector instance.
type NGramVectorFactory func(ch <-chan NGramNode) NGramVector

// ErrNGramOrderIsOutOfRange informs that the given NGrams is out of range for the given
var ErrNGramOrderIsOutOfRange = errors.New("nGrams order is out of range")

// nGramVectorBuilder implements NGramVectorBuilder interface
type nGramVectorBuilder struct {
	parents []NGramVector
	factory NGramVectorFactory
	tree    rbtree.Tree
}

// NGramNode represents tree node for the given nGram
type NGramNode struct {
	Key   Key
	Count WordCount
}

// Less tells is current elements is bigger than the other
func (n *NGramNode) Less(other rbtree.Item) bool {
	return n.Key < other.(*NGramNode).Key
}

// Key represents a NGramNode key as a composition of a NGram context and wordID
type Key uint64

// MakeKey creates uint64 key for the given pair (word, context)
func MakeKey(word WordID, context ContextOffset) Key {
	if context > maxContextOffset {
		log.Fatal(ErrContextOverflow)
	}

	return Key(utils.Pack(context, word))
}

// GetWordID returns the wordID for the given key
func (k Key) GetWordID() WordID {
	return utils.UnpackRight(uint64(k))
}

// GetContext returns the context for the given key
func (k Key) GetContext() ContextOffset {
	return utils.UnpackLeft(uint64(k))
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
			node := &NGramNode{
				Key:   MakeKey(nGram, parent),
				Count: count,
			}

			prev := m.tree.Find(node)

			if prev != nil {
				(prev.(*NGramNode)).Count += count
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
	ch := make(chan NGramNode)

	go func() {
		for iter := m.tree.NewIterator(); iter.Next() != nil; {
			node := iter.Get().(*NGramNode)
			ch <- *node
		}

		close(ch)
	}()

	return m.factory(ch)
}
