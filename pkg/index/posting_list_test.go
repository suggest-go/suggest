package index

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/suggest-go/suggest/pkg/compression"
	"github.com/suggest-go/suggest/pkg/store"
)

func TestSkipping(t *testing.T) {
	skipEncoder, _ := compression.SkippingEncoder(3)

	postings := []struct {
		name    string
		posting PostingList
		encoder compression.Encoder
	}{
		{
			name:    "skipping",
			posting: &skippingPostingList{skippingGap: 3},
			encoder: skipEncoder,
		},
		{
			name:    "dummy",
			posting: &postingList{},
			encoder: compression.VBEncoder(),
		},
		{
			name:    "bitmap",
			posting: &bitmapPostingList{},
			encoder: compression.BitmapEncoder(),
		},
	}

	cases := []struct {
		name          string
		list          []uint32
		to            uint32
		lowerBound    uint32
		tail          []uint32
		expectedError bool
	}{
		{
			name:       "#1",
			list:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:         1,
			lowerBound: 1,
			tail:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
		},
		{
			name:       "#2",
			list:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:         2,
			lowerBound: 13,
			tail:       []uint32{13, 29, 101, 506, 10003, 10004, 12000, 12001},
		},
		{
			name:       "#3",
			list:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:         12000,
			lowerBound: 12000,
			tail:       []uint32{12000, 12001},
		},
		{
			name:       "#4",
			list:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:         12001,
			lowerBound: 12001,
			tail:       []uint32{12001},
		},
		{
			name:       "#5",
			list:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:         0,
			lowerBound: 1,
			tail:       []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
		},
		{
			name:          "#6",
			list:          []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001},
			to:            12002,
			lowerBound:    0,
			tail:          []uint32{},
			expectedError: true,
		},
	}

	for _, p := range postings {
		t.Run(fmt.Sprintf("Test %s posting", p.name), func(t *testing.T) {
			posting := p.posting
			encoder := p.encoder

			for _, c := range cases {
				buf := &bytes.Buffer{}
				out := store.NewBytesOutput(buf)

				if _, err := encoder.Encode(c.list, out); err != nil {
					t.Errorf("Unexpected error occurs: %v", err)
				}

				err := posting.Init(PostingListContext{
					ListSize: len(c.list),
					Reader:   store.NewBytesInput(buf.Bytes()),
				})

				if err != nil {
					t.Errorf("Unexpected error occurs: %v", err)
				}

				actual := []uint32{}
				v, err := posting.LowerBound(c.to)

				if v != c.lowerBound {
					t.Errorf("Test %s fail, expected %v, got %v", c.name, c.lowerBound, v)
				}

				if err != nil && !c.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}

				for !c.expectedError {
					v, err := posting.Get()

					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}

					actual = append(actual, v)

					if !posting.HasNext() {
						break
					}

					v, err = posting.Next()

					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				}

				if !reflect.DeepEqual(c.tail, actual) {
					t.Errorf("Test %s fail, expected posting list: %v, got: %v", c.name, c.tail, actual)
				}
			}
		})
	}
}

func BenchmarkDummyNext(b *testing.B) {
	benchmarkNext(b, &postingList{}, compression.VBEncoder())
}

func BenchmarkSkippingNext(b *testing.B) {
	encoder, _ := compression.SkippingEncoder(64)
	benchmarkNext(b, &skippingPostingList{skippingGap: 64}, encoder)
}

func BenchmarkBitmapNext(b *testing.B) {
	benchmarkNext(b, &bitmapPostingList{}, compression.BitmapEncoder())
}

func benchmarkNext(b *testing.B, posting PostingList, encoder compression.Encoder) {
	for _, n := range []int{65, 256, 650, 6500, 65000, 650000} {
		b.Run(fmt.Sprintf("Iterate %d", n), func(b *testing.B) {
			list := make([]uint32, 0, n)

			for i := 0; i < cap(list); i++ {
				list = append(list, uint32(i))
			}

			buf := &bytes.Buffer{}
			out := store.NewBytesOutput(buf)

			if _, err := encoder.Encode(list, out); err != nil {
				b.Errorf("Unexpected error occurs: %v", err)
			}

			in := store.NewBytesInput(buf.Bytes())
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = in.Seek(0, io.SeekStart)

				err := posting.Init(PostingListContext{
					ListSize: len(list),
					Reader:   in,
				})

				if err != nil {
					b.Fatalf("Unexpected error: %v", err)
				}

				for j := uint32(1); j < uint32(n) && posting.HasNext(); j++ {
					v, err := posting.Next()

					if err != nil {
						b.Fatalf("Unexpected error: %v", err)
					}

					if j != v {
						b.Fatalf("Should receive %d, got %d", j, v)
					}
				}
			}
		})
	}
}

func BenchmarkDummyLowerBound(b *testing.B) {
	benchmarkLowerBound(b, &postingList{}, compression.VBEncoder())
}

func BenchmarkSkippingLowerBound(b *testing.B) {
	encoder, _ := compression.SkippingEncoder(64)
	benchmarkLowerBound(b, &skippingPostingList{skippingGap: 64}, encoder)
}

func BenchmarkBitmapLowerBound(b *testing.B) {
	benchmarkLowerBound(b, &bitmapPostingList{}, compression.BitmapEncoder())
}

func benchmarkLowerBound(b *testing.B, posting PostingList, encoder compression.Encoder) {
	for _, n := range []int{65, 256, 650, 6500, 65000, 650000} {
		b.Run(fmt.Sprintf("LowerBound %d", n), func(b *testing.B) {
			list := make([]uint32, 0, n)

			for i := 0; i < cap(list); i++ {
				list = append(list, uint32(i))
			}

			buf := &bytes.Buffer{}
			out := store.NewBytesOutput(buf)

			if _, err := encoder.Encode(list, out); err != nil {
				b.Errorf("Unexpected error occurs: %v", err)
			}

			in := store.NewBytesInput(buf.Bytes())
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = in.Seek(0, io.SeekStart)

				err := posting.Init(PostingListContext{
					ListSize: n,
					Reader:   in,
				})

				if err != nil {
					b.Fatalf("Unexpected error %v", err)
				}

				to := uint32(i % n)
				v, err := posting.LowerBound(to)

				if err != nil {
					b.Fatalf("Unexpected error %v", err)
				}

				if v != to {
					b.Fatalf("Test fail, expected %v, got %v", to, v)
				}
			}
		})
	}
}
