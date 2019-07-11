// +build ignore

package index

import (
	"bytes"
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
	"log"
	"testing"
)

func TestSkipping(t *testing.T) {
	list := []uint32{1, 13, 29, 101, 506, 10003, 10004, 12000, 12001}
	encoder := compression.SkippingEncoder(4)
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

	t.Error("Fail")
	/*
		for i := 0; i < 100; i++ {
			v, err := posting.Get()

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			log.Printf("Parse %v\n", v)

			if !posting.HasNext() {
				break
			}

			v, err = posting.Next()

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}
	*/

	posting.init(&postingListContext{
		listSize: len(list),
		reader:   store.NewBytesInput(buf.Bytes()),
		decoder:  nil,
	})

	v, err := posting.LowerBound(9001)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	log.Printf("LowerBound %v\n", v)

	for i := 0; i < 100; i++ {
		v, err := posting.Get()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		log.Printf("Parse %v\n", v)

		if !posting.HasNext() {
			break
		}

		v, err = posting.Next()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
