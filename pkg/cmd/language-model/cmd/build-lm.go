package cmd

import (
	"fmt"

	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildLMCmd)
}

var buildLMCmd = &cobra.Command{
	Use:   "build-lm -c [config path]",
	Short: "builds ngram language model for the given config",
	Long:  `builds ngram language model for the given config and saves it in the binary format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lm.ReadConfig(configPath)

		if err != nil {
			return fmt.Errorf("Couldn't read a config %v", err)
		}

		return lm.StoreBinaryLMFromGoogleFormat(config)
	},
}
