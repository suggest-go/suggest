package cmd

import (
	"github.com/suggest-go/suggest/internal/suggest/api"
	"log"

	"github.com/spf13/cobra"
)

var (
	port string
)

func init() {
	suggestCmd.Flags().StringVarP(&port, "port", "p", "8080", "listen port")

	rootCmd.AddCommand(suggestCmd)
}

var suggestCmd = &cobra.Command{
	Use:   "service-run -c [config path] -p [port]",
	Short: "runs http server",
	Long:  "runs http server with REST API",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("suggest: ")
		log.SetFlags(0)

		config := api.AppConfig{
			Port:       port,
			ConfigPath: configPath,
			PidPath:    pidPath,
		}

		app := api.NewApp(config)

		return app.Run()
	},
}
