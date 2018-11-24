package main

import (
	"flag"
	"github.com/alldroll/suggest/cmd/suggest/api"
)

var (
	configPath string
	port       string
)

func init() {
	flag.StringVar(&configPath, "config", "resources/index_config.json", "config path")
	flag.StringVar(&port, "port", "8080", "listen port")
}

func main() {
	config := api.AppConfig{
		Port:       port,
		ConfigPath: configPath,
	}

	app := api.NewApp(config)
	app.Run()
}
