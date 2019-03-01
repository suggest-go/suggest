package language_model

type NGramReader interface {
	Read() error
}

type googleNGramFormatReader struct {
	indexer    Indexer
	sourcePath string
}

func NewGoogleNGramReader(indexer Indexer, sourcePath string) *googleNGramFormatReader {
	return &googleNGramFormatReader{
		indexer:    indexer,
		sourcePath: sourcePath,
	}
}
