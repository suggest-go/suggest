package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(evalCmd)
}

var evalCmd = &cobra.Command{
	Use:   "eval -c [config path]",
	Short: "cli to approximate string search access",
	Long:  `cli to approximate string search access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(configPath)

		if err != nil {
			return fmt.Errorf("Failed to open config file: %v", err)
		}

		config, err := lm.ReadConfig(f)

		if err != nil {
			return fmt.Errorf("Failed to read config file: %v", err)
		}

		languageModel, err := lm.RetrieveLMFromBinary(config)

		if err != nil {
			return err
		}

		tokenizer := lm.NewTokenizer(config.GetWordsAlphabet())
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">> ")

		for scanner.Scan() {
			sentence := tokenizer.Tokenize(strings.TrimSpace(scanner.Text()))

			if len(sentence) == 0 {
				fmt.Print(">> ")
				continue
			}

			start := time.Now()
			score, err := languageModel.ScoreSentence(sentence)
			elapsed := time.Since(start).String()

			if err != nil {
				return err
			}

			fmt.Printf("Sentence: %v, Score:%f, Elapsed: %s\n", sentence, score, elapsed)
			fmt.Print(">> ")
		}

		return scanner.Err()
	},
}
