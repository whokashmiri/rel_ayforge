package config

import (
	"fmt"
	"net"
	"os"
	"time"
)

type Config struct {
	HTTPAddr            string
	HTTPReadTimeout     time.Duration
	HTTPWriteTimeout    time.Duration
	HTTPShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: getEnv("HTTP_Addr", ":8080"),
	}
	var err error
	cfg.HTTPReadTimeout, err = parseDuration(
		"HTTP_READ_TIMEOUT",
		"10s",
	)
	if err != nil {
		return Config{}, err
	}
	cfg.HTTPWriteTimeout, err = parseDuration(
		"HTTP_WRITE_TIMEOUT",
		"15s",
	)

	if err != nil {
		return Config{}, err
	}

	cfg.HTTPShutdownTimeout, err = parseDuration(
		"HTTP_SHUTDOWN_TIMEOUT",
		"10s",
	)
	if err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func parseDuration(key, fallback string) (time.Duration, error) {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)

	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater then zero", key)
	}
	return duration, nil
}

func (cfg Config) Validate() error {
	if cfg.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR cannot empty")
	}
	if _, _, err := net.SplitHostPort(cfg.HTTPAddr); err != nil {
		return fmt.Errorf("invalid HTTP_ADDR: %w", err)
	}
	if cfg.HTTPReadTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_TIMEOUT must be positive")
	}
	if cfg.HTTPWriteTimeout <= 0 {
		return fmt.Errorf("HTTP_WRITE_TIMEOUT must be positive")
	}
	if cfg.HTTPShutdownTimeout <= 0 {
		return fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT must be positive")
	}
	return nil
}
