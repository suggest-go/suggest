package api

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/suggest-go/suggest/internal/http"
	"github.com/suggest-go/suggest/internal/spellchecker/dep"
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/suggest"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// App is our application
type App struct {
	config AppConfig
}

// AppConfig is an application config
type AppConfig struct {
	Port       string
	ConfigPath string
	PidPath    string
	IndexDescription suggest.IndexDescription
}

// NewApp creates new instance of App for the given config
func NewApp(config AppConfig) App {
	return App{
		config: config,
	}
}

// Run starts the application
// performs http requests handling
func (a App) Run() error {
	config, err := lm.ReadConfig(a.config.ConfigPath)

	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	spellchecker, err := dep.BuildSpellChecker(config, a.config.IndexDescription)

	if err != nil {
		return err
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		a.listenToSystemSignals(cancelFn)
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/predict/{query}/", (&predictHandler{spellchecker}).handle).Methods("GET")

	corsHeaders := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET"})

	handler := handlers.LoggingHandler(os.Stdout, r)
	handler = handlers.CORS(corsHeaders, corsMethods)(handler)
	httpServer := http.NewServer(handler, "0.0.0.0:"+a.config.Port)

	return httpServer.Run(ctx)
}

// configureService tries to retrieve index descriptions and to setup the suggest service
func (a App) configureService(suggestService *suggest.Service) error {
	description, err := suggest.ReadConfigs(a.config.ConfigPath)

	if err != nil {
		return err
	}

	for _, config := range description {
		if err := suggestService.AddIndexByDescription(config); err != nil {
			return err
		}
	}

	return nil
}

// listenToSystemSignals handles OS signals
func (a App) listenToSystemSignals(cancelFn context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Println("Interrupt signal..")
			cancelFn()
		}
	}
}
