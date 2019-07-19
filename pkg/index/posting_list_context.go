package index

import (
	"github.com/alldroll/suggest/pkg/store"
)

// PostingListContext is the entity that holds context information
// for the corresponding Posting List
type PostingListContext interface {
	// GetListSize returns a size of the posting list
	GetListSize() int
	// GetReader returns a configured reader for the posting list
	GetReader() store.Input
}

// postingListContext implements the PostingListContext interface
type postingListContext struct {
	listSize int
	reader   store.Input
}

// GetListSize returns a size of the posting list
func (c *postingListContext) GetListSize() int {
	return c.listSize
}

// GetReader returns a configured reader for the posting list
func (c *postingListContext) GetReader() store.Input {
	return c.reader
}
