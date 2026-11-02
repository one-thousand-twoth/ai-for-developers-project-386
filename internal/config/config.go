// Package config validates the server's environment before opening resources.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
	_ "time/tzdata" // Keep IANA names available in CGO-free deployments.
)

type Config struct {
	DatabasePath string
	TimeZone     *time.Location
	ListenAddr   string
}

func Load() (Config, error) {
	name := os.Getenv("CALENDAR_TIMEZONE")
	if name == "" || name == "Local" || strings.TrimSpace(name) != name {
		return Config{}, fmt.Errorf("CALENDAR_TIMEZONE must be a named IANA time zone")
	}
	zone, err := time.LoadLocation(name)
	if err != nil {
		return Config{}, fmt.Errorf("CALENDAR_TIMEZONE must be a named IANA time zone: %w", err)
	}
	return Config{
		DatabasePath: envOr("DATABASE_PATH", "./data/calendar.sqlite"),
		TimeZone:     zone,
		ListenAddr:   envOr("LISTEN_ADDR", ":8080"),
	}, nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
