// Package mph represents minimal perfect hash function implementation
package mph

import (
	"fmt"
	"math"
	"sort"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/store"
)

// MPH represents minimal perfect hash function
type MPH interface {
	// Build builds a MPH for the given dictionary
	Build(dict dictionary.Dictionary) error
	// Get returns a hash value for the given word
	Get(word dictionary.Value) dictionary.Key
	// Store stores the given MPH structure into output
	Store(out store.Output) (int, error)
	// Load loads from the input a MPH structure
	Load(in store.Input) (int, error)
}

// New creates a new instance of MPH object
func New() MPH {
	return &mph{
		auxiliary: []int32{},
		values:    []uint32{},
	}
}

// mph implements MPH interface
type mph struct {
	auxiliary []int32
	values    []dictionary.Key
}

// Build builds a MPH for the given dictionary
// Inspired by http://stevehanov.ca/blog/?id=119
func (m *mph) Build(dict dictionary.Dictionary) error {
	size := uint32(dict.Size()) //cdb dictionary can not be more that uint32 size
	buckets := make([][]dictionary.Key, size)
	auxiliary := make([]int32, size, size)
	values := make([]dictionary.Key, 0, size)

	// Step 1: Place all of the keys into buckets
	err := dict.Iterate(func(key dictionary.Key, value dictionary.Value) error {
		d := hash(0, value) % size
		buckets[d] = append(buckets[d], key)
		values = append(values, math.MaxUint32)

		return nil
	})

	if err != nil {
		return err
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
		d := uint32(1)
		slots := make([]uint32, 0, len(bucket))

		// Repeatedly try different values of d until we find a hash function
		// that places all items in the bucket into free slots
		for item < len(bucket) {
			value, err := dict.Get(bucket[item])

			if err != nil {
				return fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
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
			return fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
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
	freeSlots := make([]int, 0, size)

	for i, val := range values {
		if val == math.MaxUint32 {
			freeSlots = append(freeSlots, i)
		}
	}

	for _, bucket := range buckets[bucketIter:] {
		if len(bucket) == 0 || len(freeSlots) == 0 {
			break
		}

		slot := freeSlots[len(freeSlots)-1]
		freeSlots = freeSlots[:len(freeSlots)-1]
		val, err := dict.Get(bucket[0])

		if err != nil {
			return fmt.Errorf("Failed to get bucket's key from the dictionary: %v", err)
		}

		// We subtract one to ensure it's negative even if the zeroeth slot was
		// used.
		auxiliary[hash(0, val)%size] = int32(-slot - 1)
		values[slot] = bucket[0]
	}

	m.auxiliary = auxiliary
	m.values = values

	return nil
}

// Get returns a hash value for the given word
func (m *mph) Get(word dictionary.Value) dictionary.Key {
	d := m.auxiliary[hash(0, word)%uint32(len(m.auxiliary))]

	if d < 0 {
		return m.values[-d-1]
	}

	return m.values[hash(uint32(d), word)%uint32(len(m.values))]
}

// Store stores the given MPH structure into output
func (m *mph) Store(out store.Output) (int, error) {
	n, err := out.WriteUInt32(uint32(len(m.values)))

	if err != nil {
		return n, fmt.Errorf("failed to write the length of values: %v", err)
	}

	for _, v := range m.values {
		s, err := out.WriteUInt32(v)
		n += s

		if err != nil {
			return n, fmt.Errorf("failed to write a value: %v", err)
		}
	}

	s, err := out.WriteUInt32(uint32(len(m.auxiliary)))
	n += s

	if err != nil {
		return n, fmt.Errorf("failed to write the length of auxiliary: %v", err)
	}

	for _, v := range m.auxiliary {
		s, err := out.WriteUInt32(uint32(v))
		n += s

		if err != nil {
			return n, fmt.Errorf("failed to writer a value: %v", err)
		}
	}

	return n, nil
}

// Load loads from the input a MPH structure
func (m *mph) Load(in store.Input) (int, error) {
	n, err := in.ReadUInt32()

	if err != nil {
		return 0, fmt.Errorf("failed to read the length of values: %v", err)
	}

	m.values = make([]dictionary.Key, n)

	for i := range m.values {
		v, err := in.ReadUInt32()

		if err != nil {
			return 0, fmt.Errorf("failed to read a value: %v", err)
		}

		m.values[i] = v
	}

	s, err := in.ReadUInt32()

	if err != nil {
		return 0, fmt.Errorf("failed to read the length of auxiliary: %v", err)
	}

	m.auxiliary = make([]int32, s)

	for i := range m.auxiliary {
		v, err := in.ReadUInt32()

		if err != nil {
			return 0, fmt.Errorf("failed to read a value: %v", err)
		}

		m.auxiliary[i] = int32(v)
	}

	return int(n+s)*4 + 8, nil
}

// hash encodes the given value and the salt
func hash(h uint32, value string) uint32 {
	if h == 0 {
		h = 2166136261
	}

	for _, c := range []byte(value) {
		h *= 16777619
		h ^= uint32(c)
	}

	return h
}

// has tells is the given value in the slice
func has(slice []uint32, value uint32) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
