package main

import (
	"fmt"
	"os"
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
	return nil
}
