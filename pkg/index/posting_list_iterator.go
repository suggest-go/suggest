package index

import (
	"fmt"
	"github.com/alldroll/suggest/pkg/merger"
)

// postingListDecodedBufSize is a buffer size for raw decoded posting list
const postingListDecodedBufSize = 512

// postingListIterator is a dummy implementation of merger.ListIterator
type postingListIterator struct {
	*merger.SliceIterator
	buf []Position
}

// init initialize the iterator by the given PostingList context
func (i *postingListIterator) init(context PostingListContext) error {
	var postingList []Position

	if context.GetListSize() > postingListDecodedBufSize {
		postingList = make([]Position, context.GetListSize())
	} else {
		if i.buf == nil {
			i.buf = make([]Position, postingListDecodedBufSize)
		}

		postingList = i.buf
	}

	n, err := context.GetDecoder().Decode(context.GetReader(), postingList)

	if err != nil {
		return fmt.Errorf("failed to decode posting list: %v", err)
	}

	if i.SliceIterator != nil {
		i.SliceIterator.Reset(postingList[:n])
	} else {
		i.SliceIterator = merger.NewSliceIterator(postingList[:n])
	}

	return nil
}
