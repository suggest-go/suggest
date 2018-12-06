package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/alldroll/suggest/pkg/alphabet"
	lm "github.com/alldroll/suggest/pkg/language_model"
	"log"
	"os"
	"strings"
)

const (
	startSymbol = "<S>"
	endSymbol   = "</S>"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "config file")
}

//
func buildNGramsCount(config *lm.Config, indexer lm.Indexer) lm.Trie {
	sourceFile, err := os.Open(config.SourcePath)
	if err != nil {
		log.Fatalf("could read source file %s", err)
	}

	defer sourceFile.Close()

	retriever := lm.NewSentenceRetriver(
		lm.NewTokenizer(alphabet.CreateAlphabet(config.Alphabet)),
		bufio.NewReader(sourceFile),
		alphabet.CreateAlphabet(config.Separators),
	)

	nGramOrder := uint8(config.NGramOrder)
	trie := lm.NewTrie()
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
			s := make([]lm.WordId, 0)

			for _, nGrams := range nGramsSet {
				s = s[:0]

				for _, nGram := range nGrams {
					s = append(s, indexer.GetOrCreate(nGram))
				}

				trie.Put(s)
			}
		}
	}

	return trie
}

//
func storeNGramsCount(config *lm.Config, trie lm.Trie, indexer lm.Indexer) {
	bufs := []*bufio.Writer{}

	for i := 0; i < int(config.NGramOrder); i++ {
		f, err := os.OpenFile(fmt.Sprintf("%s/%d-gm", config.OutputPath, i+1), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		bufs = append(bufs, buf)
		defer buf.Flush()
	}

	grams := lm.NGram{}

	trie.Walk(func(path []lm.WordId, count lm.WordCount) {
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

func main() {
	log.SetPrefix("indexer: ")
	log.SetFlags(0)

	flag.Parse()

	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("could not open config file %s", err)
	}

	defer configFile.Close()

	config, err := lm.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("could read config %s", err)
	}

	indexer := lm.NewIndexer()
	trie := buildNGramsCount(config, indexer)
	storeNGramsCount(config, trie, indexer)
}
