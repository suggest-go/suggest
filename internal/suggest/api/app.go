package api

import (
	"context"
	"fmt"
	"github.com/suggest-go/suggest/internal/http"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/suggest-go/suggest/pkg/suggest"
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
	if err := a.writePIDFile(); err != nil {
		return err
	}

	suggestService := suggest.NewService()
	reindexJob := func() error {
		return a.configureService(suggestService)
	}

	if err := reindexJob(); err != nil {
		return fmt.Errorf("Fail to configure service: %w", err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		a.listenToSystemSignals(
			cancelFn,
			func() {
				if err := reindexJob(); err != nil {
					log.Printf("Fail to reload index %s", err)
				} else {
					log.Printf("Reindex done!")
				}
			},
		)
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/autocomplete/{dict}/{query}/", (&autocompleteHandler{suggestService}).handle).Methods("GET")
	r.HandleFunc("/suggest/{dict}/{query}/", (&suggestHandler{suggestService}).handle).Methods("GET")
	r.HandleFunc("/dict/list/", (&dictionaryHandler{suggestService}).handle).Methods("GET")
	r.HandleFunc("/internal/reindex/", (&reindexHandler{reindexJob}).handle).Methods("POST")

	corsHeaders := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET"})

	handler := handlers.LoggingHandler(os.Stdout, r)
	handler = handlers.CORS(corsHeaders, corsMethods)(handler)
	httpServer := http.NewServer(handler, "0.0.0.0:"+a.config.Port)

	return httpServer.Run(ctx)
}

// writePIDFile performs writing a PID of the application service
func (a App) writePIDFile() error {
	if a.config.PidPath == "" {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(a.config.PidPath), 0700)
	if err != nil {
		return fmt.Errorf("There is no such file for %s, %w", a.config.PidPath, err)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(a.config.PidPath, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Fail to write pid file to %s, %w", a.config.PidPath, err)
	}

	return nil
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
func (a App) listenToSystemSignals(cancelFn context.CancelFunc, reindexFn func()) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			log.Println("Reindex job..")
			reindexFn()

		case <-signalChan:
			log.Println("Interrupt signal..")
			cancelFn()
		}
	}
}
