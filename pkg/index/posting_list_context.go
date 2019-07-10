package index

import (
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
)

// PostingListContext is the entity that holds context information
// for the corresponding Posting List
type PostingListContext interface {
	// GetListSize returns a size of the posting list
	GetListSize() int
	// GetReader returns a configured reader for the posting list
	GetReader() store.Input
	// GetDecoder returns a decoder which should be used to decode the posting list
	GetDecoder() compression.Decoder
}

// postingListContext implements the PostingListContext interface
type postingListContext struct {
	listSize int
	reader   store.Input
	decoder  compression.Decoder
}

// GetListSize returns a size of the posting list
func (c *postingListContext) GetListSize() int {
	return c.listSize
}

// GetReader returns a configured reader for the posting list
func (c *postingListContext) GetReader() store.Input {
	return c.reader
}

// GetDecoder returns a decoder which should be used to decode the posting list
func (c *postingListContext) GetDecoder() compression.Decoder {
	return c.decoder
}
