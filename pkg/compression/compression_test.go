package compression

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/store"
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

	testCases := []struct {
		p []uint32
	}{
		{[]uint32{824, 829, 215406}},
		{[]uint32{1, 9, 13, 180, 999, 12345}},
		{[]uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12901}},
	}

	for _, ins := range instances {
		encoder := ins.encoder
		decoder := ins.decoder

		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("%s_%d", ins.name, i+1), func(t *testing.T) {
				buf := &bytes.Buffer{}
				output := store.NewBytesOutput(buf)
				list := make([]uint32, len(testCase.p))

				_, err := encoder.Encode(testCase.p, output)
				assert.NoError(t, err)

				in := store.NewBytesInput(buf.Bytes())
				_, err = decoder.Decode(in, list)

				assert.NoError(t, err)
				assert.Equal(t, list, testCase.p)
			})
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
	output := store.NewBytesOutput(buf)

	if _, err := encoder.Encode(list, output); err != nil {
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
