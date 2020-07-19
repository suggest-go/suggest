package lm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"sort"
	"strconv"

	"github.com/alldroll/go-datastructures/rbtree"
)

type packedArray struct {
	keys       []ContextOffset
	containers []arrayContainer
	values     []WordCount
	total      WordCount
}

type arrayContainer struct {
	words   []WordID
	offsets ContextOffset
}

// CreatePackedArray creates a NGramVector form the given Tree
func CreatePackedArray(tree rbtree.Tree) NGramVector {
	var node *nGramNode
	values := make([]WordCount, 0, tree.Len())
	total := WordCount(0)

	var container *arrayContainer
	context := InvalidContextOffset
	keys := []ContextOffset{}
	containers := []arrayContainer{}

	for i, iter := 0, tree.NewIterator(); iter.Next() != nil; i++ {
		node = iter.Get().(*nGramNode)
		values = append(values, node.value)
		total += node.value
		currContext := getContext(node.key)

		if container == nil || context != currContext {
			if container != nil {
				containers = append(containers, *container)
			}

			container = &arrayContainer{
				words:   []WordID{},
				offsets: ContextOffset(i),
			}

			keys = append(keys, currContext)
			context = currContext
		}

		container.words = append(container.words, getWordID(node.key))
	}

	if container != nil {
		containers = append(containers, *container)
	}

	return &packedArray{
		keys:       keys,
		containers: containers,
		values:     values,
		total:      total,
	}
}

// GetCount returns WordCount and Node ContextOffset for the given pair (word, context)
func (s *packedArray) GetCount(word WordID, context ContextOffset) (WordCount, ContextOffset) {
	i := s.find(word, context)

	if InvalidContextOffset == i {
		return 0, InvalidContextOffset
	}

	return s.values[int(i)], i
}

// GetContextOffset returns the given node context offset
func (s *packedArray) GetContextOffset(word WordID, context ContextOffset) ContextOffset {
	return s.find(word, context)
}

// CorpusCount returns size of all counts in the collection
func (s *packedArray) CorpusCount() WordCount {
	return s.total
}

// SubVector returns NGramVector for the given context
func (s *packedArray) SubVector(context ContextOffset) NGramVector {
	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= context })

	if i < 0 || i >= len(s.keys) || s.keys[i] != context {
		return nil
	}

	j := s.containers[i].offsets

	container := arrayContainer{
		words:   s.containers[i].words,
		offsets: 0,
	}

	return &packedArray{
		keys:       []ContextOffset{context},
		containers: []arrayContainer{container},
		values:     s.values[j : j+ContextOffset(len(container.words))],
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
	for i, context := range s.keys {
		for _, word := range s.containers[i].words {
			el := makeKey(word, context)
			keyEndPos += binary.PutUvarint(encodedKeys[keyEndPos:], el-prevKey)
			prevKey = el
		}
	}

	valEndPos := len(s.values) * 4
	encodedValues := make([]byte, valEndPos)

	for i, el := range s.values {
		binary.LittleEndian.PutUint32(encodedValues[i*4:(i+1)*4], el)
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

	var container *arrayContainer
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
				s.containers = append(s.containers, *container)
			}

			container = &arrayContainer{
				words:   []WordID{},
				offsets: ContextOffset(i),
			}

			s.keys = append(s.keys, currContext)
			context = currContext
		}

		container.words = append(container.words, getWordID(prev))
		keyEndPos += n
	}

	if container != nil {
		s.containers = append(s.containers, *container)
	}

	encodedValues := buf.Next(valSize)
	s.values = make([]WordCount, valSize/4)

	for i := 0; i < len(s.values); i++ {
		s.values[i] = binary.LittleEndian.Uint32(encodedValues[i*4 : (i+1)*4])
	}

	return nil
}

// find finds the given key in the collection. Returns ContextOffset if the key exists, otherwise returns InvalidContextOffset
func (s *packedArray) find(wordID WordID, context ContextOffset) ContextOffset {
	if len(s.keys) == 0 || s.keys[0] > context || s.keys[len(s.keys)-1] < context {
		return InvalidContextOffset
	}

	i := sort.Search(len(s.keys), func(i int) bool { return s.keys[i] >= context })

	if i < 0 || i >= len(s.keys) || s.keys[i] != context {
		return InvalidContextOffset
	}

	container := s.containers[i]

	j := sort.Search(len(container.words), func(i int) bool { return container.words[i] >= wordID })

	if j < 0 || j >= len(container.words) || container.words[j] != wordID {
		return InvalidContextOffset
	}

	return ContextOffset(j) + container.offsets
}

func init() {
	gob.Register(&packedArray{})
}
