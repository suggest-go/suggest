package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alldroll/suggest/pkg/dictionary"
	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/alldroll/suggest/pkg/spellchecker"
	"github.com/alldroll/suggest/pkg/suggest"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "TODO", "TODO")
}

func main() {
	flag.Parse()

	config, err := lm.ReadConfig(configPath)

	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	languageModel, err := lm.RetrieveLMFromBinary(config)

	if err != nil {
		log.Fatal(err)
	}

	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		log.Fatal(err)
	}

	tokenizer := lm.NewTokenizer(config.GetWordsAlphabet())

	// describe index configuration
	indexDescription := suggest.IndexDescription{
		Driver:    suggest.RAMDriver,
		Name:      "words",
		NGramSize: 2,
		Wrap:      [2]string{"^", "$"},
		Pad:       "$",
		Alphabet:  []string{"english", "russian", "numbers", "$^"},
	}

	// create runtime search index builder
	builder, err := suggest.NewRAMBuilder(dict, indexDescription)

	if err != nil {
		log.Fatal(err)
	}

	index, err := builder.Build()

	if err != nil {
		log.Fatal(err)
	}

	service := spellchecker.New(
		index,
		languageModel,
		tokenizer,
		dict,
	)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")

	for scanner.Scan() {
		sentence := scanner.Text()

		if len(sentence) == 0 {
			fmt.Print(">> ")
			continue
		}

		start := time.Now()
		result, err := service.Predict(sentence, 5)
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
