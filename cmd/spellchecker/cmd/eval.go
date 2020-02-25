package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/suggest-go/suggest/internal/spellchecker/dep"
	"os"
	"time"

	"github.com/suggest-go/suggest/pkg/lm"
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
			return fmt.Errorf("failed to read config file: %w", err)
		}

		service, err := dep.BuildSpellChecker(config, indexDescription)

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
