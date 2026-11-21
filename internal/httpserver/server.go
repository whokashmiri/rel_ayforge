package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/whokashmiri/rel_ayforge/internal/config"
	"github.com/whokashmiri/rel_ayforge/internal/handler"
	"github.com/whokashmiri/rel_ayforge/internal/middleware"
)

func New(cfg config.Config, logger *slog.Logger, tenantHandler *handler.TenantHandler) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging(logger))

	router.Get("/health", handler.Health)
	router.Get("/ready", handler.Ready)

	router.Post("/api/v1/tenants", tenantHandler.Create)
	router.Get("/api/v1/tenants/{id}", tenantHandler.GetByID)
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
