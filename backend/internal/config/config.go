// Package config membaca konfigurasi dari environment variable.
package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Env           string // development · staging · production
	HTTPAddr      string
	PublicBaseURL string
	DatabaseURL   string
	DefaultLang   string
	LogLevel      string
}

func Load() (Config, error) {
	c := Config{
		Env:           get("APP_ENV", "development"),
		HTTPAddr:      get("HTTP_ADDR", ":8080"),
		PublicBaseURL: strings.TrimRight(get("PUBLIC_BASE_URL", "http://localhost:8080"), "/"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		DefaultLang:   get("DEFAULT_LANG", "en"),
		LogLevel:      get("LOG_LEVEL", "info"),
	}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL belum diisi (lihat .env.example)")
	}
	switch c.DefaultLang {
	case "en", "km", "zh":
	default:
		return c, fmt.Errorf("DEFAULT_LANG harus en, km, atau zh")
	}
	return c, nil
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
