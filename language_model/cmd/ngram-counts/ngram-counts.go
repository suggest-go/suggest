package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/alldroll/suggest/alphabet"
	lm "github.com/alldroll/suggest/language_model"
	"log"
	"os"
	"strings"
)

const (
	startSymbol = "<S>"
	endSymbol   = "</S>"
)

var (
	sourcePath string
)

func init() {
	flag.StringVar(&sourcePath, "source", "", "source path")
}

func main() {
	log.SetPrefix("indexer: ")
	log.SetFlags(0)

	flag.Parse()

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatalf("could not open source file %s", err)
	}

	defer sourceFile.Close()

	retriever := lm.NewSentenceRetriver(
		lm.NewTokenizer(alphabet.NewEnglishAlphabet()),
		bufio.NewReader(sourceFile),
		alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'}),
	)

	nGramOrder := uint8(3)
	trie := lm.NewTrie()
	indexer := lm.NewIndexer()
	generators := []lm.Generator{}

	for i := uint8(1); i <= nGramOrder; i++ {
		generators = append(
			generators,
			lm.NewGenerator(
				i,
				startSymbol,
				endSymbol,
			),
		)
	}

	for {
		sentence := retriever.Retrieve()
		if sentence == nil {
			break
		}

		if len(sentence) == 0 {
			continue
		}

		for i := 0; i < int(nGramOrder); i++ {
			generator := generators[i]
			nGramsSet := generator.Generate(sentence)
			s := make([]uint32, 0)

			for _, nGrams := range nGramsSet {
				s = s[:0]

				for _, nGram := range nGrams {
					s = append(s, indexer.GetOrCreate(nGram))
				}

				trie.Put(s)
			}
		}
	}

	bufs := []*bufio.Writer{}

	for i := 0; i < int(nGramOrder); i++ {
		f, err := os.OpenFile(fmt.Sprintf("%d-gm", i+1), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		bufs = append(bufs, buf)
		defer buf.Flush()
	}

	grams := lm.NGram{}

	trie.Walk(func(path []uint32, count uint32) {
		if len(path) == 0 {
			return
		}

		buf := bufs[len(path)-1]
		grams = grams[:0]

		for _, g := range path {
			nGram, err := indexer.Find(g)
			if err != nil {
				panic(err)
			}

			grams = append(grams, nGram)
		}

		fmt.Fprintf(buf, "%s\t%d\n", strings.Join(grams, " "), count)
	})
}
