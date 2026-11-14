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

	// File upload (sementara disk lokal; nanti Cloudflare R2)
	StorageDir string

	// Terjemahan AI
	Translator          string // claude · google · fake · off
	AnthropicKey        string
	AnthropicModel      string
	GoogleTranslateKey  string
	TranslationFallback string // google · none
	WorkerInServe       bool   // true = worker (terjemahan) ikut jalan di `phsar serve`

	// Google Maps / Places di browser (Seller Center, admin). Kunci browser: batasi per domain di Google Cloud.
	GoogleMapsKey   string
	GoogleMapsMapID string
}

func Load() (Config, error) {
	loadDotEnv(".env", "../.env")
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
		StorageDir:          get("STORAGE_DIR", "./data/uploads"),

		Translator:          get("TRANSLATOR", "claude"),
		AnthropicKey:        os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicModel:      get("ANTHROPIC_MODEL", "claude-haiku-5-5"),
		GoogleTranslateKey:  os.Getenv("GOOGLE_TRANSLATE_API_KEY"),
		TranslationFallback: get("TRANSLATION_FALLBACK", "google"),
		WorkerInServe:       get("WORKER_IN_SERVE", "true") == "true",
		GoogleMapsKey:       os.Getenv("GOOGLE_MAPS_API_KEY"),
		GoogleMapsMapID:     get("GOOGLE_MAPS_MAP_ID", "DEMO_MAP_ID"),
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
	switch c.Translator {
	case "claude", "google", "off":
	case "fake":
		if c.Env == "production" {
			return c, fmt.Errorf("TRANSLATOR=fake tidak boleh dipakai di production")
		}
	default:
		return c, fmt.Errorf("TRANSLATOR harus claude, google, fake, atau off")
	}
	switch c.DefaultLang {
	case "en", "km", "zh":
	default:
		return c, fmt.Errorf("DEFAULT_LANG harus en, km, atau zh")
	}
	return c, nil
}

// loadDotEnv membaca file .env (folder sekarang atau folder induk) supaya `go run ./cmd/phsar ...`
// langsung jalan tanpa make. Variabel yang sudah ada di environment tidak ditimpa.
func loadDotEnv(paths ...string) {
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, v, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			k, v = strings.TrimSpace(strings.TrimPrefix(k, "export ")), strings.TrimSpace(v)
			v = strings.Trim(v, `"'`)
			if _, exists := os.LookupEnv(k); !exists {
				os.Setenv(k, v)
			}
		}
		return
	}
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
