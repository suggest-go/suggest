package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"

	"github.com/spf13/cobra"
)

func init() {
	evalCmd.Flags().StringVarP(&dict, "dict", "d", "", "dictionary name")
	evalCmd.MarkPersistentFlagRequired("dict")

	rootCmd.AddCommand(evalCmd)
}

var evalCmd = &cobra.Command{
	Use:   "eval -c [config path] -d [dict]",
	Short: "cli to approximate string search access",
	Long:  `cli to approximate string search access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("eval: ")
		log.SetFlags(0)

		suggestService, err := configureService()

		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print(">> ")

		for scanner.Scan() {
			searchConf, err := suggest.NewSearchConfig(scanner.Text(), 5, metric.CosineMetric(), 0.4)

			if err != nil {
				return err
			}

			result, err := suggestService.Suggest(dict, searchConf)

			if err != nil {
				return err
			}

			fmt.Println()

			for _, item := range result {
				fmt.Printf("%s, score: %f\n", item.Value, item.Score)
			}

			fmt.Print(">> ")
		}

		return scanner.Err()
	},
}

// configureService creates and configures suggest service for the given
// config and the dictionary
func configureService() (*suggest.Service, error) {
	config, err := os.Open(configPath)

	if err != nil {
		return nil, fmt.Errorf("Failed to open config file: %v", err)
	}

	defer config.Close()
	description, err := suggest.ReadConfigs(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to read configs: %v", err)
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
