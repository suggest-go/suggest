package suggest

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	instances := []struct {
		encoder Encoder
		decoder Decoder
	}{
		{BinaryEncoder(), BinaryDecoder()},
		{VBEncoder(), VBDecoder()},
	}

	cases := []struct {
		p PostingList
	}{
		{PostingList{824, 829, 215406}},
		{PostingList{1, 9, 13, 180, 999, 12345}},
	}

	for _, ins := range instances {
		encoder := ins.encoder
		decoder := ins.decoder

		for _, c := range cases {
			bytes := encoder.Encode(c.p)
			list := decoder.Decode(bytes)

			if !reflect.DeepEqual(list, c.p) {
				t.Errorf("Fail, expected posting list: %v, got: %v", c.p, list)
			}
		}
	}
}

func BenchmarkBinaryDecode(b *testing.B) {
	benchmarkDecode(BinaryEncoder(), BinaryDecoder(), b)
}

func BenchmarkVBDecode(b *testing.B) {
	benchmarkDecode(VBEncoder(), VBDecoder(), b)
}

func benchmarkDecode(encoder Encoder, decoder Decoder, b *testing.B) {
	list := make(PostingList, 0, 1000)

	for i := 1; i <= 1000; i++ {
		list = append(list, Position(rand.Intn(10000)))
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	bytes := encoder.Encode(list)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder.Decode(bytes)
	}
}
