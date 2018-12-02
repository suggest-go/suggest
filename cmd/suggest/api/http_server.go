package api

import (
	"context"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	r    http.Handler
	addr string
}

//
func newHttpServer(r http.Handler, addr string) *httpServer {
	return &httpServer{
		r:    r,
		addr: addr,
	}
}

//
func (h *httpServer) Run(ctx context.Context) error {
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
