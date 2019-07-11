package index

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
)

func TestSkipping(t *testing.T) {
	list := []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001}
	encoder := compression.SkippingEncoder()
	buf := &bytes.Buffer{}

	if _, err := encoder.Encode(list, buf); err != nil {
		t.Errorf("Unexpected error occurs: %v", err)
	}

	posting := &skippingPostingList{}

	posting.init(&postingListContext{
		listSize: len(list),
		reader:   store.NewBytesInput(buf.Bytes()),
		decoder:  nil,
	})

	actual := []uint32{}

	for {
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

	if !reflect.DeepEqual(list, actual) {
		t.Errorf("Test fail, expected posting list: %v, got: %v", list, actual)
	}

	actual = actual[:0]

	posting.init(&postingListContext{
		listSize: len(list),
		reader:   store.NewBytesInput(buf.Bytes()),
		decoder:  nil,
	})

	v, err := posting.LowerBound(400)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if v != 506 {
		t.Errorf("Test fail, expected 506, got %v", v)
	}

	for {
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

	if !reflect.DeepEqual(list[4:], actual) {
		t.Errorf("Test fail, expected posting list: %v, got: %v", list[4:], actual)
	}
}
