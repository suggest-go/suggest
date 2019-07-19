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
	skipEnc, _ := SkippingEncoder(3)
	skipDec, _ := SkippingDecoder(3)

	instances := []struct {
		name    string
		encoder Encoder
		decoder Decoder
	}{
		{"binary", BinaryEncoder(), BinaryDecoder()},
		{"varint", VBEncoder(), VBDecoder()},
		{"skipping", skipEnc, skipDec},
	}

	cases := []struct {
		p []uint32
	}{
		{[]uint32{824, 829, 215406}},
		{[]uint32{1, 9, 13, 180, 999, 12345}},
		{[]uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12901}},
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

func BenchmarkSkippingDecode(b *testing.B) {
	enc, err := SkippingEncoder(64)

	if err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	dec, err := SkippingDecoder(64)

	if err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	benchmarkDecode(enc, dec, b)
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

	if _, err := encoder.Encode(list, buf); err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	encoded := buf.Bytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		in := store.NewBytesInput(encoded)

		if _, err := decoder.Decode(in, list); err != nil {
			b.Errorf("Unexpected error: %v", err)
		}

		buf.Reset()
	}
}
