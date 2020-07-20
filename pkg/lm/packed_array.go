package lm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"sort"
	"strconv"

	"github.com/alldroll/go-datastructures/rbtree"
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

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (s *packedArray) MarshalBinary() ([]byte, error) {
	var result bytes.Buffer

	encodedKeys := make([]byte, len(s.values)*binary.MaxVarintLen64)
	prevKey := uint64(0)
	keyEndPos := 0

	// performs delta encoding
	for _, container := range s.containers {
		for i := container.from; i < container.to; i++ {
			row := s.values[i]
			el := makeKey(utils.UnpackLeft(row), container.context)
			keyEndPos += binary.PutUvarint(encodedKeys[keyEndPos:], el-prevKey)
			prevKey = el
		}
	}

	valEndPos := len(s.values) * 4
	encodedValues := make([]byte, valEndPos)

	for i, el := range s.values {
		binary.LittleEndian.PutUint32(encodedValues[i*4:(i+1)*4], utils.UnpackRight(el))
	}

	// allocate buffer capacity
	result.Grow(keyEndPos + valEndPos + 4 + strconv.IntSize*2)

	// write header
	if _, err := fmt.Fprintln(&result, keyEndPos, valEndPos, s.total); err != nil {
		return nil, err
	}

	// write data
	result.Write(encodedKeys[:keyEndPos])
	result.Write(encodedValues)

	return result.Bytes(), nil
}

// UnmarshalBinary decodes the binary form
func (s *packedArray) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	keySize, valSize := 0, 0

	if _, err := fmt.Fscanln(buf, &keySize, &valSize, &s.total); err != nil {
		return err
	}

	var container *rangeContainer
	keyEndPos := 0
	encodedKeys := buf.Next(keySize)
	prev := uint64(0)
	context := InvalidContextOffset

	// 0, 1, 3, 6, 7 -> 0, 1, 4, 10, 17 (delta decoding)
	for i := 0; i < valSize/4; i++ {
		data, n := binary.Uvarint(encodedKeys[keyEndPos:])
		data += prev
		prev = data
		currContext := getContext(prev)

		if container == nil || context != currContext {
			if container != nil {
				container.to = uint32(i)
				s.containers = append(s.containers, *container)
			}

			container = &rangeContainer{
				context: currContext,
				from:    uint32(i),
			}

			context = currContext
		}

		s.values = append(s.values, utils.Pack(getWordID(data), 0)) // fill it later
		keyEndPos += n
	}

	if container != nil {
		container.to = uint32(len(s.values))
		s.containers = append(s.containers, *container)
	}

	encodedValues := buf.Next(valSize)

	for i := 0; i < len(s.values); i++ {
		count := binary.LittleEndian.Uint32(encodedValues[i*4 : (i+1)*4])
		s.values[i] = utils.Pack(utils.UnpackLeft(s.values[i]), count)
	}

	return nil
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

	return values[j], ContextOffset(j)
}

func init() {
	gob.Register(&packedArray{})
}
