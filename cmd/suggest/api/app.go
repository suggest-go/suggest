package api

import (
	"context"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
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
	PidPath    string
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

	if err := a.writePIDFile(); err != nil {
		log.Fatal(err)
	}

	suggestService := suggest.NewService()
	if err := a.configureService(suggestService); err != nil {
		log.Fatal(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		a.listenToSystemSignals(
			cancelFn,
			func() {
				if err := a.configureService(suggestService); err != nil {
					// TODO use proper way to say about errors
					log.Fatal(err)
				}
			},
		)
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/autocomplete/{dict}/{query}/", (&autocompleteHandler{suggestService}).handle)
	r.HandleFunc("/suggest/{dict}/{query}/", (&suggestHandler{suggestService}).handle)
	r.HandleFunc("/dict/list/", (&dictionaryHandler{suggestService}).handle)

	httpServer := newHttpServer(r, "0.0.0.0:"+port)
	httpServer.Run(ctx)
}

//
func (a App) writePIDFile() error {
	if a.config.PidPath == "" {
		log.Fatal("PidPath should be set")
	}

	err := os.MkdirAll(filepath.Dir(a.config.PidPath), 0700)
	if err != nil {
		log.Fatal("There is no such file for %s, %s", a.config.PidPath, err)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(a.config.PidPath, []byte(pid), 0644); err != nil {
		// TODO make corresponding error
		return err
	}

	return nil
}

//
func (a App) configureService(suggestService *suggest.Service) error {
	config, err := os.Open(a.config.ConfigPath)
	if err != nil {
		return err
	}

	defer config.Close()
	description, err := suggest.ReadConfigs(config)

	if err != nil {
		return err
	}

	for _, config := range description {
		if err := suggestService.AddOnDiscIndex(config); err != nil {
			return err
		}
	}

	return nil
}

//
func (a App) listenToSystemSignals(cancelFn context.CancelFunc, reindexFn func()) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			log.Println("signupSignal")
			reindexFn()

		case <-signalChan:
			log.Println("signtermSignal")
			cancelFn()
		}
	}
}
