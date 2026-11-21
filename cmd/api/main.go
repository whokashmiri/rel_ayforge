package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/whokashmiri/rel_ayforge/internal/config"
	"github.com/whokashmiri/rel_ayforge/internal/httpserver"
)

func main() {
	fmt.Println("Starting RelayForce API...")
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "RelayForce failed", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("RelayForge API initialized successfully")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	server := httpserver.New(cfg, logger)
	fmt.Println("RelayForge listening on", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("start HTTP server : %w", err)
		}
	}
	return nil
}
