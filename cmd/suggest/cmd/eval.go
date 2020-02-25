package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"

	"github.com/spf13/cobra"
)

var (
	topK       int
	similarity float64
)

func init() {
	evalCmd.Flags().StringVarP(&dict, "dict", "d", "", "dictionary name")
	evalCmd.MarkPersistentFlagRequired("dict")

	evalCmd.Flags().IntVarP(&topK, "topK", "k", 5, "topK elements")
	evalCmd.Flags().Float64VarP(&similarity, "sim", "s", 0.5, "similarity of candidates")

	rootCmd.AddCommand(evalCmd)
}

var evalCmd = &cobra.Command{
	Use:   "eval -c [config path] -d [dict]",
	Short: "cli to approximate string search access",
	Long:  `cli to approximate string search access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		suggestService, err := configureService()

		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">> ")

		for scanner.Scan() {
			query := strings.TrimSpace(scanner.Text())

			if len(query) == 0 {
				fmt.Print(">> ")
				continue
			}

			searchConf, err := suggest.NewSearchConfig(query, topK, metric.CosineMetric(), similarity)

			if err != nil {
				return err
			}

			start := time.Now()
			result, err := suggestService.Suggest(dict, searchConf)
			elapsed := time.Since(start).String()

			if err != nil {
				return err
			}

			for _, item := range result {
				fmt.Printf("%s, score: %f\n", item.Value, item.Score)
			}

			fmt.Printf("\nElapsed: %s (%d candidates)\n", elapsed, len(result))
			fmt.Print(">> ")
		}

		return scanner.Err()
	},
}

// configureService creates and configures suggest service for the given
// config and the dictionary
func configureService() (*suggest.Service, error) {
	description, err := suggest.ReadConfigs(configPath)

	if err != nil {
		return nil, fmt.Errorf("Failed to read configs: %w", err)
	}

	suggestService := suggest.NewService()

	for _, config := range description {
		if config.Name != dict {
			continue
		}

		if err := suggestService.AddIndexByDescription(config); err != nil {
			return nil, err
		}
	}

	if len(suggestService.GetDictionaries()) != 1 {
		return nil, fmt.Errorf("Dictionary %s is not found", dict)
	}

	return suggestService, nil
}
