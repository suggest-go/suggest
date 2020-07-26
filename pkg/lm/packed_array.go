package lm

import (
	"fmt"
	"io"
	"sort"

	"github.com/suggest-go/suggest/pkg/store"
	"github.com/suggest-go/suggest/pkg/utils"
)

type packedArray struct {
	containers []rangeContainer
	values     []packedValue
	total      WordCount
}

type packedValue = uint64

func newPackedValue(wordID WordID, wordCount WordCount) packedValue {
	return utils.Pack(wordID, wordCount)
}

func getWordID(p packedValue) WordID {
	return utils.UnpackLeft(p)
}

func getWordCount(p packedValue) WordCount {
	return utils.UnpackRight(p)
}

type rangeContainer = uint64

func newRangeContainer(context ContextOffset, from uint32) rangeContainer {
	return utils.Pack(context, from)
}

func getContext(rc rangeContainer) ContextOffset {
	return utils.UnpackLeft(rc)
}

func getFrom(rc rangeContainer) uint32 {
	return utils.UnpackRight(rc)
}

// NewNGramVector creates a new instance of NGramVector.
func NewNGramVector() NGramVector {
	return &packedArray{}
}

// GetCount returns WordCount and Node ContextOffset for the given pair (word, context)
func (s *packedArray) GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset) {
	value, i := s.find(word, context)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return getWordCount(value), i
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
	i := s.findContainerPos(context)

	if i == -1 {
		return nil
	}

	containers := []rangeContainer{s.containers[i]}

	if i < len(s.containers)-1 {
		containers = append(containers, s.containers[i+1])
	}

	return &packedArray{
		containers: containers,
		values:     s.values,
		total:      s.total,
	}
}

func (s *packedArray) Store(out store.Output) (int, error) {
	// write header
	containerSize := uint32(8 * len(s.containers))
	valuesSize := uint32(8 * len(s.values))
	n, err := fmt.Fprintln(out, containerSize, valuesSize, s.total)
	p := n

	if err != nil {
		return p, err
	}

	// write values
	if n, err = out.Write(uint64SliceAsByteSlice(s.containers)); err != nil {
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

	s.containers = byteSliceAsUint64Slice(data[:containersSize])
	s.values = byteSliceAsUint64Slice(data[containersSize:])

	return p, nil
}

// find finds the given key in the collection. Returns ContextOffset if the key exists, otherwise returns InvalidContextOffset
func (s *packedArray) find(wordID WordID, context ContextOffset) (packedValue, ContextOffset) {
	i := s.findContainerPos(context)

	if i == -1 {
		return 0, InvalidContextOffset
	}

	values := []packedValue{}
	container := s.containers[i]
	from := getFrom(container)

	if i == len(s.containers)-1 {
		values = s.values[from:]
	} else {
		to := getFrom(s.containers[i+1])
		values = s.values[from:to]
	}

	target := newPackedValue(wordID, 0)
	j := sort.Search(len(values), func(i int) bool { return values[i] >= target })

	if j < 0 || j >= len(values) || getWordID(values[j]) != wordID {
		return 0, InvalidContextOffset
	}

	return values[j], from + ContextOffset(j)
}

// findContainerPos searches for a position of a range container with the given context.
// Returns -1 if container does not exist.
func (s *packedArray) findContainerPos(context ContextOffset) int {
	if len(s.containers) == 0 || getContext(s.containers[0]) > context || getContext(s.containers[len(s.containers)-1]) < context {
		return -1
	}

	target := newRangeContainer(context, 0)
	i := sort.Search(len(s.containers), func(i int) bool { return s.containers[i] >= target })

	if i < 0 || i >= len(s.containers) || getContext(s.containers[i]) != context {
		return -1
	}

	return i
}

// CreatePackedArray creates a NGramVector form the channel of NGramNodes.
func CreatePackedArray(ch <-chan NGramNode) NGramVector {
	pa := &packedArray{
		containers: []rangeContainer{},
		values:     []packedValue{},
		total:      0,
	}

	context := InvalidContextOffset
	from := uint32(0)

	for node := range ch {
		pa.total += node.Count
		currContext := node.Key.GetContext()

		if context != currContext || len(pa.containers) == 0 {
			pa.containers = append(pa.containers, newRangeContainer(currContext, from))
			context = currContext
		}

		pa.values = append(pa.values, newPackedValue(node.Key.GetWordID(), node.Count))
		from++
	}

	return pa
}
