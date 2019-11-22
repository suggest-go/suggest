package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	r    http.Handler
	addr string
}

// NewServer creates new instance of HttpServer
func NewServer(r http.Handler, addr string) *Server {
	return &Server{
		r:    r,
		addr: addr,
	}
}

// Run starts serving http requests
func (h *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         h.addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.r,
	}

	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Fail to shutdown server %s", err)
		}
	}()

	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Println("Server was shutdown gracefully")
		return nil
	}

	return err
}
