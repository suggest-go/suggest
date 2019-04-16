package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alldroll/suggest/pkg/alphabet"
	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

func init() {
	countNGramsCmd.Flags().StringVarP(&configPath, "config-file", "c", "", "TODO describe usage")
	countNGramsCmd.MarkFlagRequired("config-file")

	rootCmd.AddCommand(countNGramsCmd)
}

var countNGramsCmd = &cobra.Command{
	Use:   "ngram-count -c [config path]",
	Short: "builds ngram counts for the given config file using google ngram format",
	Long:  `builds ngram counts for the given config file using google ngram format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, err := os.Open(configPath)
		if err != nil {
			return fmt.Errorf("could not open config file %s", err)
		}

		defer configFile.Close()

		config, err := lm.ReadConfig(configFile)

		if err != nil {
			return fmt.Errorf("could read config %s", err)
		}

		trie, err := buildNGramsCount(config)

		if err != nil {
			return err
		}

		return storeNGramsCount(config, trie)
	},
}

// buildNGramsCount builds a count trie
func buildNGramsCount(config *lm.Config) (lm.CountTrie, error) {
	sourceFile, err := os.Open(config.SourcePath)

	if err != nil {
		return nil, fmt.Errorf("could read source file %s", err)
	}

	defer sourceFile.Close()

	retriever := lm.NewSentenceRetriever(
		lm.NewTokenizer(alphabet.CreateAlphabet(config.Alphabet)),
		bufio.NewReader(sourceFile),
		alphabet.CreateAlphabet(config.Separators),
	)

	builder := lm.NewNGramBuilder(
		config.StartSymbol,
		config.EndSymbol,
	)

	return builder.Build(retriever, config.NGramOrder), nil
}

// storeNGramsCount flushes the constructed count trie on FS
func storeNGramsCount(config *lm.Config, trie lm.CountTrie) error {
	writer := lm.NewGoogleNGramWriter(config.NGramOrder, config.OutputPath)

	if err := writer.Write(trie); err != nil {
		return fmt.Errorf("could save ngrams %s", err)
	}

	return nil
}
