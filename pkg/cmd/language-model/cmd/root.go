package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lm",
	Short: "cli to interact with a ngram language model",
	Long:  `cli to interact with a ngram languauge model`,
}

// Execute runs commands handling
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
