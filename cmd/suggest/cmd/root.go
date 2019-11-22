package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	pidPath    string
)

var rootCmd = &cobra.Command{
	Use:   "suggest",
	Short: "cli to interact with a suggest index and http server",
	Long:  `cli to interact with a suggest index and http server`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to the config file")
	rootCmd.MarkPersistentFlagFilename("config")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.PersistentFlags().StringVarP(&pidPath, "pid", "", "", "path to pid file")
}

// Execute runs commands handling
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
