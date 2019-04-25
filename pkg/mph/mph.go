// Package mph represents minimal perfect hash function implementation
package mph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"sort"

	"github.com/alldroll/suggest/pkg/dictionary"
)

// MPH represents minimal perfect hash function
type MPH interface {
	// Get returns a hash value for the given word
	Get(word dictionary.Value) dictionary.Key
}

// BuildMPH builds a MPH for the given dictionary
// Inspired by http://stevehanov.ca/blog/?id=119
func BuildMPH(dict dictionary.Dictionary) (MPH, error) {
	size := dict.Size()

	buckets := make([][]dictionary.Key, size, size)
	auxiliary := make([]int32, size, size)
	values := make([]dictionary.Key, 0, size)

	// Step 1: Place all of the keys into buckets
	err := dict.Iterate(func(key dictionary.Key, value dictionary.Value) error {
		d := int(hash(0, value)) % size
		buckets[d] = append(buckets[d], key)
		values = append(values, math.MaxUint32)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Step 2: Sort the buckets and process the ones with the most items first.
	sort.Slice(buckets, func(i, j int) bool {
		return len(buckets[i]) >= len(buckets[j])
	})

	bucketIter := 0

	for _, bucket := range buckets {
		if len(bucket) <= 1 {
			break
		}

		item := 0
		d := int32(1)
		slots := make([]int, 0, len(bucket))

		// Repeatedly try different values of d until we find a hash function
		// that places all items in the bucket into free slots
		for item < len(bucket) {
			value, err := dict.Get(bucket[item])

			if err != nil {
				return nil, fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
			}

			slot := hash(d, value) % size

			if values[slot] != math.MaxUint32 || has(slots, slot) {
				d++
				item = 0
				slots = slots[:0]
			} else {
				slots = append(slots, slot)
				item++
			}
		}

		val, err := dict.Get(bucket[0])

		if err != nil {
			return nil, fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
		}

		auxiliary[hash(0, val)%size] = int32(d)

		for i, key := range bucket {
			values[slots[i]] = key
		}

		bucketIter++
	}

	// Only buckets with 1 item remain. Process them more quickly by directly
	// placing them into a free slot. Use a negative value of d to indicate
	// this.
	freeslots := make([]int, 0, size)
	for i, val := range values {
		if val == math.MaxUint32 {
			freeslots = append(freeslots, i)
		}
	}

	for _, bucket := range buckets[bucketIter:] {
		if len(bucket) == 0 || len(freeslots) == 0 {
			break
		}

		slot := freeslots[len(freeslots)-1]
		freeslots = freeslots[:len(freeslots)-1]
		val, err := dict.Get(bucket[0])

		if err != nil {
			return nil, fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
		}

		// We subtract one to ensure it's negative even if the zeroeth slot was
		// used.
		auxiliary[hash(0, val)%size] = int32(-slot - 1)
		values[slot] = bucket[0]
	}

	return &mph{
		auxiliary: auxiliary,
		values:    values,
	}, nil
}

// mph implements MPH interface
type mph struct {
	auxiliary []int32
	values    []dictionary.Key
}

// Get returns a hash value for the given word
func (m *mph) Get(word dictionary.Value) dictionary.Key {
	d := m.auxiliary[hash(0, word)%len(m.auxiliary)]

	if d < 0 {
		return m.values[-d-1]
	}

	return m.values[hash(d, word)%len(m.values)]
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (m *mph) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(m.auxiliary); err != nil {
		return nil, err
	}

	if err := encoder.Encode(m.values); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the binary data
func (m *mph) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	if err := decoder.Decode(&m.auxiliary); err != nil {
		return err
	}

	return decoder.Decode(&m.values)
}

// hash encodes the given value and the salt
func hash(d int32, value string) int {
	h := uint64(d)

	if h == 0 {
		h = 0x811C9DC5
	}

	for _, c := range []byte(value) {
		h = (h ^ uint64(c)*16777619) & 0xffffffff
	}

	return int(h)
}

// has tells is the given value in the slice
func has(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

func init() {
	gob.Register(&mph{})
}
