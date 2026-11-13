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

	// Login
	JWTSecret           string // minimal 32 karakter
	TelegramBotToken    string // dari @BotFather; kosong = login Telegram nonaktif
	TelegramBotUsername string // tanpa @, mis. phsar_login_bot
	TelegramMode        string // polling (lokal) · webhook (server) · off
	TelegramWebhookKey  string // bagian rahasia di URL /webhooks/telegram/<key>
	SMSProvider         string // log (pengembangan: kode OTP ditulis ke log) · nanti provider sungguhan
}

func Load() (Config, error) {
	c := Config{
		Env:           get("APP_ENV", "development"),
		HTTPAddr:      get("HTTP_ADDR", ":8080"),
		PublicBaseURL: strings.TrimRight(get("PUBLIC_BASE_URL", "http://localhost:8080"), "/"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		DefaultLang:   get("DEFAULT_LANG", "en"),
		LogLevel:      get("LOG_LEVEL", "info"),

		JWTSecret:           os.Getenv("JWT_SECRET"),
		TelegramBotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramBotUsername: strings.TrimPrefix(os.Getenv("TELEGRAM_BOT_USERNAME"), "@"),
		TelegramMode:        get("TELEGRAM_MODE", "polling"),
		TelegramWebhookKey:  os.Getenv("TELEGRAM_WEBHOOK_KEY"),
		SMSProvider:         get("SMS_PROVIDER", "log"),
	}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL belum diisi (lihat .env.example)")
	}
	if len(c.JWTSecret) < 32 {
		if c.Env != "development" {
			return c, fmt.Errorf("JWT_SECRET wajib diisi (minimal 32 karakter)")
		}
		c.JWTSecret = "dev-only-insecure-secret-change-me-0123456789"
	}
	if c.TelegramBotToken == "" {
		c.TelegramMode = "off"
	} else if c.TelegramBotUsername == "" {
		return c, fmt.Errorf("TELEGRAM_BOT_USERNAME wajib diisi kalau TELEGRAM_BOT_TOKEN diisi")
	}
	if c.Env == "production" && c.SMSProvider == "log" {
		return c, fmt.Errorf("SMS_PROVIDER=log tidak boleh dipakai di production")
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
