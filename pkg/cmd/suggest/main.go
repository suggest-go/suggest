package main

import (
	"flag"
	"github.com/alldroll/suggest/pkg/cmd/suggest/api"
	"log"
)

var (
	configPath string
	port       string
	pidPath    string
)

func init() {
	flag.StringVar(&configPath, "config", "resources/index_config.json", "path to indexes config file")
	flag.StringVar(&pidPath, "pid", "", "path to pid file")
	flag.StringVar(&port, "port", "8080", "listen port")
}

func main() {
	log.SetPrefix("suggest: ")
	log.SetFlags(0)
	flag.Parse()

	config := api.AppConfig{
		Port:       port,
		ConfigPath: configPath,
		PidPath:    pidPath,
	}

	app := api.NewApp(config)
	app.Run()
}
