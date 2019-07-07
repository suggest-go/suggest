package index

import (
	"github.com/alldroll/suggest/pkg/merger"
)

//
type PostingListIterator interface {
	merger.ListIterator
}

//
func NewPostingListIterator(slice []Position) PostingListIterator {
	return merger.NewSliceIterator(slice)
}
