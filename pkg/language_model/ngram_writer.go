package language_model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	fileFormat  = "%s/%d-gm"
	nGramFormat = "%s\t%d\n"
)

type NGramWriter interface {
	Write(trie Trie) error
}

func NewGoogleNGramWriter(indexer Indexer, nGramOrder uint8, outputPath string) *googleNGramFormatWriter {
	return &googleNGramFormatWriter{
		indexer:    indexer,
		nGramOrder: nGramOrder,
		outputPath: outputPath,
	}
}

type googleNGramFormatWriter struct {
	indexer    Indexer
	outputPath string
	nGramOrder uint8
}

func (gw *googleNGramFormatWriter) Write(trie Trie) error {
	bufs := []*bufio.Writer{}

	for i := 0; i < int(gw.nGramOrder); i++ {
		f, err := os.OpenFile(fmt.Sprintf(fileFormat, gw.outputPath, i+1), os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			return err
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		bufs = append(bufs, buf)
		defer buf.Flush()
	}

	grams := make(NGram, 0, int(gw.nGramOrder))

	trie.Walk(func(path []WordId, count WordCount) {
		if len(path) == 0 {
			return
		}

		grams = grams[:0]

		for _, g := range path {
			nGram, err := gw.indexer.Find(g)
			if err != nil {
				panic(err) // TODO catch error
			}

			grams = append(grams, nGram)
		}

		fmt.Fprintf(bufs[len(path)-1], nGramFormat, strings.Join(grams, " "), count)
	})

	return nil
}
