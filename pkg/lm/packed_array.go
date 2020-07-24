package lm

import (
	"fmt"
	"io"
	"sort"

	"github.com/alldroll/go-datastructures/rbtree"
	"github.com/suggest-go/suggest/pkg/store"
	"github.com/suggest-go/suggest/pkg/utils"
)

type packedArray struct {
	containers []rangeContainer
	values     []packedValue // we can store it on disc
	total      WordCount
}

type packedValue = uint64

type rangeContainer struct {
	context  ContextOffset
	from, to uint32
}

// CreatePackedArray creates a NGramVector form the given Tree
func CreatePackedArray(tree rbtree.Tree) NGramVector {
	var node *nGramNode
	total := WordCount(0)

	var container *rangeContainer
	context := InvalidContextOffset
	values := make([]packedValue, 0, tree.Len())
	containers := []rangeContainer{}

	for i, iter := 0, tree.NewIterator(); iter.Next() != nil; i++ {
		node = iter.Get().(*nGramNode)
		total += node.value
		currContext := getContext(node.key)

		if container == nil || context != currContext {
			if container != nil {
				container.to = uint32(i)
				containers = append(containers, *container)
			}

			container = &rangeContainer{
				context: currContext,
				from:    uint32(i),
			}

			context = currContext
		}

		values = append(values, utils.Pack(getWordID(node.key), node.value))
	}

	if container != nil {
		container.to = uint32(len(values))
		containers = append(containers, *container)
	}

	return &packedArray{
		containers: containers,
		values:     values,
		total:      total,
	}
}

// GetCount returns WordCount and Node ContextOffset for the given pair (word, context)
func (s *packedArray) GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset) {
	value, i := s.find(word, context)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return utils.UnpackRight(value), i
}

// GetContextOffset returns the given node context offset
func (s *packedArray) GetContextOffset(word WordID, context ContextOffset) ContextOffset {
	_, i := s.find(word, context)

	return i
}

// CorpusCount returns size of all counts in the collection
func (s *packedArray) CorpusCount() WordCount {
	return s.total
}

// SubVector returns NGramVector for the given context
func (s *packedArray) SubVector(context ContextOffset) NGramVector {
	if len(s.containers) == 0 || s.containers[0].context > context || s.containers[len(s.containers)-1].context < context {
		return nil
	}

	i := sort.Search(len(s.containers), func(i int) bool { return s.containers[i].context >= context })

	if i < 0 || i >= len(s.containers) || s.containers[i].context != context {
		return nil
	}

	return &packedArray{
		containers: []rangeContainer{s.containers[i]},
		values:     s.values,
		total:      s.total,
	}
}

func (s *packedArray) Store(out store.Output) (int, error) {
	// write header
	containerSize := uint32(12 * len(s.containers))
	valuesSize := uint32(8 * len(s.values))
	n, err := fmt.Fprintln(out, containerSize, valuesSize, s.total)
	p := n

	if err != nil {
		return p, err
	}

	// write values
	if n, err = out.Write(rangeContainerSliceAsByteSlice(s.containers)); err != nil {
		return p + n, err
	}

	p += n

	if n, err = out.Write(uint64SliceAsByteSlice(s.values)); err != nil {
		return p + n, err
	}

	return p + n, nil
}

func (s *packedArray) Load(in store.Input) (int, error) {
	containersSize, valuesSize := uint32(0), uint32(0)
	n, err := fmt.Fscanln(in, &containersSize, &valuesSize, &s.total)

	if err != nil {
		return 0, err
	}

	offset, err := in.Seek(0, io.SeekCurrent)
	p := int(offset)

	if err != nil {
		return p, err
	}

	var data []byte

	if accessable, ok := in.(store.SliceAccessible); ok {
		n := int(containersSize + valuesSize)
		data = accessable.Data()[p : p+n]

		if _, err = in.Seek(offset+int64(n), io.SeekStart); err != nil {
			return p, err
		}

		p += n
	} else {
		data = make([]byte, containersSize+valuesSize)

		if n, err = in.Read(data); err != nil {
			return p + n, err
		}

		p += n
	}

	s.containers = byteSliceAsRangeContainerSlice(data[:containersSize])
	s.values = byteSliceAsUint64Slice(data[containersSize:])

	return p, nil
}

// find finds the given key in the collection. Returns ContextOffset if the key exists, otherwise returns InvalidContextOffset
func (s *packedArray) find(wordID WordID, context ContextOffset) (packedValue, ContextOffset) {
	if len(s.containers) == 0 || s.containers[0].context > context || s.containers[len(s.containers)-1].context < context {
		return 0, InvalidContextOffset
	}

	i := sort.Search(len(s.containers), func(i int) bool { return s.containers[i].context >= context })

	if i < 0 || i >= len(s.containers) || s.containers[i].context != context {
		return 0, InvalidContextOffset
	}

	container := s.containers[i]
	values := s.values[container.from:container.to]
	target := utils.Pack(wordID, 0)
	j := sort.Search(len(values), func(i int) bool { return values[i] >= target })

	if j < 0 || j >= len(values) || utils.UnpackLeft(values[j]) != wordID {
		return 0, InvalidContextOffset
	}

	return values[j], container.from + ContextOffset(j)
}
