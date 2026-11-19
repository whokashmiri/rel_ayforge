package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/whokashmiri/rel_ayforge/internal/config"
	"github.com/whokashmiri/rel_ayforge/internal/handler"
)

func New(cfg config.Config) *http.Server {
	router := chi.NewRouter()

	router.Get("/health", handler.Health)
	router.Get("/ready", handler.Ready)
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       60 * time.Second,
	}
	return server

}
