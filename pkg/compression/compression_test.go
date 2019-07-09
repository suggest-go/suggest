package compression

import (
	"bytes"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/alldroll/suggest/pkg/store"
)

func TestEncodeDecode(t *testing.T) {
	instances := []struct {
		name    string
		encoder Encoder
		decoder Decoder
	}{
		{"binary", BinaryEncoder(), BinaryDecoder()},
		{"varint", VBEncoder(), VBDecoder()},
	}

	cases := []struct {
		p []uint32
	}{
		{[]uint32{824, 829, 215406}},
		{[]uint32{1, 9, 13, 180, 999, 12345}},
	}

	for _, ins := range instances {
		encoder := ins.encoder
		decoder := ins.decoder

		for _, c := range cases {
			buf := &bytes.Buffer{}
			list := make([]uint32, len(c.p))

			if _, err := encoder.Encode(c.p, buf); err != nil {
				t.Errorf("[%s] Unexpected error occurs: %v", ins.name, err)
			}

			in := store.NewBytesInput(buf.Bytes())

			if _, err := decoder.Decode(in, list); err != nil {
				t.Errorf("[%s] Unexpected error occurs: %v", ins.name, err)
			}

			if !reflect.DeepEqual(list, c.p) {
				t.Errorf("Fail [%s], expected posting list: %v, got: %v", ins.name, c.p, list)
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
	list := make([]uint32, 0, 1000)

	for i := 1; i <= 1000; i++ {
		list = append(list, uint32(rand.Intn(10000)))
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	buf := &bytes.Buffer{}
	encoder.Encode(list, buf)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		in := store.NewBytesInput(buf.Bytes())
		decoder.Decode(in, list)
		buf.Reset()
	}
}
