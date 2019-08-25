package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alldroll/suggest/pkg/store"

	"github.com/alldroll/suggest/pkg/dictionary"
	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/alldroll/suggest/pkg/spellchecker"
	"github.com/alldroll/suggest/pkg/suggest"
)

var (
	configPath string

	indexDescription = suggest.IndexDescription{
		Driver:    suggest.RAMDriver,
		Name:      "words",
		NGramSize: 3,
		Wrap:      [2]string{"^", "$"},
		Pad:       "$",
		Alphabet:  []string{"english", "russian", "numbers", "$^'"},
	}

	topK = 5

	similarity = 0.5
)

func init() {
	flag.StringVar(&configPath, "config", "spellchecker config", "spellchecker configuration file")
}

func main() {
	flag.Parse()
	config, err := lm.ReadConfig(configPath)

	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	service, err := buildSpellChecker(config)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")

	for scanner.Scan() {
		sentence := scanner.Text()

		if len(sentence) == 0 {
			fmt.Print(">> ")
			continue
		}

		start := time.Now()
		result, err := service.Predict(sentence, topK, similarity)
		elapsed := time.Since(start).String()

		if err != nil {
			log.Fatal(err)
		}

		for _, item := range result {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\nElapsed: %s (%d candidates)\n", elapsed, len(result))
		fmt.Print(">> ")
	}

	log.Fatal(scanner.Err())
}

// buildSpellChecker builds spellchecker
func buildSpellChecker(config *lm.Config) (*spellchecker.SpellChecker, error) {
	directory, err := store.NewFSDirectory(config.GetOutputPath())

	if err != nil {
		return nil, fmt.Errorf("failed to create a fs directory: %v", err)
	}

	languageModel, err := lm.RetrieveLMFromBinary(directory, config)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a lm model from binary format: %v", err)
	}

	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		return nil, fmt.Errorf("failed to open a cdb dictionary: %v", err)
	}

	// create runtime search index builder
	builder, err := suggest.NewRAMBuilder(dict, indexDescription)

	if err != nil {
		return nil, fmt.Errorf("failed to create a ngram index: %v", err)
	}

	index, err := builder.Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build a ngram index: %v", err)
	}

	return spellchecker.New(
		index,
		languageModel,
		lm.NewTokenizer(config.GetWordsAlphabet()),
		dict,
	), nil
}
