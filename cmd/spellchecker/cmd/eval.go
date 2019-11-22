package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"

	"github.com/suggest-go/suggest/pkg/store"

	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/spellchecker"
	"github.com/suggest-go/suggest/pkg/suggest"
)

var (
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
	rootCmd.AddCommand(evalCmd)
}

var evalCmd = &cobra.Command{
	Use:   "eval -c [config path]",
	Short: "spellchecker cli",
	Long:  `spellchecker cli`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lm.ReadConfig(configPath)

		if err != nil {
			return fmt.Errorf("failed to read config file: %v", err)
		}

		service, err := buildSpellChecker(config)

		if err != nil {
			return err
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
				return err
			}

			for _, item := range result {
				fmt.Printf("%s\n", item)
			}

			fmt.Printf("\nElapsed: %s (%d candidates)\n", elapsed, len(result))
			fmt.Print(">> ")
		}

		return scanner.Err()
	},
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
