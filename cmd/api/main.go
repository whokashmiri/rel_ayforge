package main

import (
	"fmt"
	"os"

	"github.com/whokashmiri/rel_ayforge/internal/config"
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
	fmt.Println("Configuration loaded successfully")
	fmt.Println("HTTP address:", cfg.HTTPAddr)
	fmt.Println("Read timeout:", cfg.HTTPReadTimeout)
	fmt.Println("Write timeout:", cfg.HTTPWriteTimeout)
	fmt.Println("Shutdown timeout:", cfg.HTTPShutdownTimeout)

	fmt.Println("RelayForge API initialized successfully")
	return nil
}
