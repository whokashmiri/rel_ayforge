package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":8080")
	t.Setenv("HTTP_READ_TIMEOUT", "10s")
	t.Setenv("HTTP_WRITE_TIMEOUT", "15s")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "10s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Errorf(
			"HTTPAddr = %q, want %q",
			cfg.HTTPAddr,
			":8080",
		)
	}

	if cfg.HTTPReadTimeout != 10*time.Second {
		t.Errorf(
			"HTTPReadTimeout = %v, want 10s",
			cfg.HTTPReadTimeout,
		)
	}

	if cfg.HTTPWriteTimeout != 15*time.Second {
		t.Errorf(
			"HTTPWriteTimeout = %v, want 15s",
			cfg.HTTPWriteTimeout,
		)
	}

	if cfg.HTTPShutdownTimeout != 10*time.Second {
		t.Errorf(
			"HTTPShutdownTimeout = %v, want 10s",
			cfg.HTTPShutdownTimeout,
		)
	}
}

func TestLoadInvalidDuration(t *testing.T) {
	t.Setenv("HTTP_READ_TIMEOUT", "invalid")

	_, err := Load()

	if err == nil {
		t.Fatal("Load() expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "HTTP_READ_TIMEOUT") {
		t.Errorf(
			"error = %q, want HTTP_READ_TIMEOUT",
			err,
		)
	}
}

func TestLoadInvalidAddress(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")

	_, err := Load()

	if err == nil {
		t.Fatal("Load() expected an error for empty address")
	}
}
