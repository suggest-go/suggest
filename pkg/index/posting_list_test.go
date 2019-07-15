package index

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
)

func TestSkipping(t *testing.T) {
	cases := []struct {
		name       string
		list       []uint32
		to         uint32
		lowerBound uint32
		tail       []uint32
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
	}

	posting := &skippingPostingList{}

	for _, c := range cases {
		encoder, _ := compression.SkippingEncoder(3)
		buf := &bytes.Buffer{}

		if _, err := encoder.Encode(c.list, buf); err != nil {
			t.Errorf("Unexpected error occurs: %v", err)
		}

		posting.init(&postingListContext{
			listSize: len(c.list),
			reader:   store.NewBytesInput(buf.Bytes()),
			decoder:  nil,
		})

		actual := []uint32{}
		v, err := posting.LowerBound(c.to)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if v != c.lowerBound {
			t.Errorf("Test %s fail, expected %v, got %v", c.name, c.lowerBound, v)
		}

		for i := 0; i < 30; i++ {
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
}
