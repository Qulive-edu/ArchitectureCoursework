package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/config"

	"github.com/go-chi/chi/v5"
)

type HttpServer struct {
	server *http.Server
	notify chan error
}

func NewServer(cfg config.Server, handler http.Handler) *HttpServer {
	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.Timeout * 2,
	}
	h := &HttpServer{server: srv, notify: make(chan error, 1)}
	go func() {
		h.notify <- srv.ListenAndServe()
		close(h.notify)
	}()
	return h
}

func (h *HttpServer) Notify() <-chan error { return h.notify }

func (h *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return h.server.Shutdown(ctx)
}

// helper to mount handlers
func mount(r chi.Router, pattern string, h http.HandlerFunc) {
	r.Method("GET", pattern, h)
	r.Method("POST", pattern, h)
}
