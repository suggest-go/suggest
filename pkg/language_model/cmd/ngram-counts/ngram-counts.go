package main

import (
	"bufio"
	"flag"
	"github.com/alldroll/suggest/pkg/alphabet"
	lm "github.com/alldroll/suggest/pkg/language_model"
	"log"
	"os"
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

	builder := lm.NewNGramBuilder(retriever, indexer, config.NGramOrder)

	return builder.Build()
}

//
func storeNGramsCount(config *lm.Config, trie lm.Trie, indexer lm.Indexer) {
	writer := lm.NewGoogleNGramWriter(indexer, config.NGramOrder, config.OutputPath)
	if err := writer.Write(trie); err != nil {
		log.Fatalf("could save ngrams %s", err)
	}
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
