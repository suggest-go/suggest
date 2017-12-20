package suggest

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestBinaryDecodeNumber(t *testing.T) {
	cases := []struct {
		p PostingList
	}{
		{PostingList{824, 829, 215406}},
		{PostingList{1, 9, 13, 180, 999, 12345}},
	}

	encoder := BinaryEncoder()
	decoder := BinaryDecoder()

	for _, c := range cases {
		bytes := encoder.Encode(c.p)
		list := decoder.Decode(bytes)

		if !reflect.DeepEqual(list, c.p) {
			t.Errorf("Fail, expected posting list: %v, got: %v", c.p, list)
		}
	}
}

func BenchmarkBytesAlgoDecode(b *testing.B) {
	encoder := BinaryEncoder()
	decoder := BinaryDecoder()

	list := make(PostingList, 0, 1000)
	for i := 1; i <= 1000; i++ {
		list = append(list, Position(rand.Intn(1000)))
	}

	bytes := encoder.Encode(list)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder.Decode(bytes)
	}
}

func BenchmarkVBDecode(b *testing.B) {
	encoder := VBEncoder()
	decoder := VBDecoder()

	list := make(PostingList, 0, 1000)
	for i := 1; i <= 1000; i++ {
		list = append(list, Position(rand.Intn(1000)))
	}

	bytes := encoder.Encode(list)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder.Decode(bytes)
	}
}
