package store

import "io"

// Output is a wrap for method Write
type Output interface {
	io.Writer
}
