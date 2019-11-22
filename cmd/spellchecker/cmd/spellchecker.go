package cmd

import (
	"github.com/suggest-go/suggest/internal/spellchecker/api"
	"log"

	"github.com/spf13/cobra"
)

var (
	port string
)

func init() {
	spellcheckerCmd.Flags().StringVarP(&port, "port", "p", "8080", "listen port")

	rootCmd.AddCommand(spellcheckerCmd)
}

var spellcheckerCmd = &cobra.Command{
	Use:   "service-run -c [config path] -p [port]",
	Short: "runs http server",
	Long:  "runs http server with REST API",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("spellchecker: ")
		log.SetFlags(0)

		config := api.AppConfig{
			Port:       port,
			ConfigPath: configPath,
			IndexDescription: indexDescription,
		}

		app := api.NewApp(config)

		return app.Run()
	},
}
