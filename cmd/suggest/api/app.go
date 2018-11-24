package api

import (
	"context"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//
type App struct {
	config AppConfig
}

//
type AppConfig struct {
	Port       string
	ConfigPath string
}

//
func NewApp(config AppConfig) App {
	return App{
		config: config,
	}
}

//
func (a App) Run() {
	port := a.config.Port
	if port == "" {
		log.Fatal("Port must be set")
	}

	service, err := a.configureService()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		a.listenToSystemSignals(cancelFn)
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/autocomplete/{dict}/{query}/", (&autocompleteHandler{service}).handle)
	r.HandleFunc("/suggest/{dict}/{query}/", (&suggestHandler{service}).handle)
	r.HandleFunc("/dict/list/", (&dictionaryHandler{service}).handle)

	httpServer := newHttpServer(r, "0.0.0.0:"+port)
	httpServer.Run(ctx)
}

//
func (a App) configureService() (*suggest.Service, error) {
	config, err := os.Open(a.config.ConfigPath)
	if err != nil {
		return nil, err
	}

	defer config.Close()

	suggestService := suggest.NewService()
	configs, err := suggest.ReadConfigs(config)

	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if err := suggestService.AddOnDiscIndex(config); err != nil {
			return nil, err
		}
	}

	return suggestService, nil
}

//
func (a App) listenToSystemSignals(cancelFn context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			log.Println("signupSignal")
			// TODO index update
		case <-signalChan:
			log.Println("signtermSignal")
			cancelFn()
		}
	}
}
