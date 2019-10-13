package index

import (
	"github.com/suggest-go/suggest/pkg/store"
)

// PostingListContext is the entity that holds context information
// for the corresponding Posting List
type PostingListContext struct {
	ListSize int
	Reader   store.Input
}

