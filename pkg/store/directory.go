package store

// Directory is a flat list of files.
// Inspired by org.apache.lucene.store.Directory
type Directory interface {
	// CreateOutput creates a new writer in the given directory with the given name
	CreateOutput(name string) (Output, error)
	// OpenInput returns a reader for the given name
	OpenInput(name string) (Input, error)
}
