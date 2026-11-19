// Package config membaca konfigurasi dari environment variable.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"phsar/internal/secretbox"
)

type Config struct {
	Env string // development · staging · production
	// AdminMFADevSkip: khusus APP_ENV=development — login admin tidak meminta kode authenticator (MFA).
	// Otomatis nyala di development; matikan dengan ADMIN_MFA_SKIP=false. Di staging/production selalu mati.
	AdminMFADevSkip bool
	HTTPAddr        string
	PublicBaseURL   string
	// Alamat Seller Center & admin kalau dipisah ke subdomain (deploy: seller.<domain>, admin.<domain>).
	// Kosong = PUBLIC_BASE_URL. Dipakai untuk link di notifikasi Telegram & push.
	SellerBaseURL string
	AdminBaseURL  string
	DatabaseURL   string
	DefaultLang   string
	LogLevel      string

	// Login
	JWTSecret           string // minimal 32 karakter
	TelegramBotToken    string // dari @BotFather; kosong = login Telegram nonaktif
	PushContact         string // PUSH_CONTACT: kontak VAPID untuk layanan Web Push (mailto:… / https://…)
	TelegramBotUsername string // tanpa @, mis. phsar_login_bot
	TelegramMode        string // polling (lokal) · webhook (server) · off
	TelegramWebhookKey  string // bagian rahasia di URL /webhooks/telegram/<key>
	SMSProvider         string // log (pengembangan: kode OTP ditulis ke log) · off (login SMS mati, hanya Telegram) · nanti provider sungguhan

	// File upload (sementara disk lokal; nanti Cloudflare R2)
	StorageDir string

	// Terjemahan AI
	Translator          string // claude · google · fake · off
	AnthropicKey        string
	AnthropicModel      string
	AnthropicBaseURL    string // askbot: kosong = https://api.anthropic.com (bisa diganti untuk tes/proxy)
	GeminiKey           string // GEMINI_API_KEY (cadangan; utamanya di Admin → Integrasi)
	GeminiModel         string
	GoogleTranslateKey  string
	TranslationFallback string // google · none
	WorkerInServe       bool   // true = worker (terjemahan) ikut jalan di `phsar serve`

	// Google Maps / Places di browser (Seller Center, admin). Kunci browser: batasi per domain di Google Cloud.
	GoogleMapsKey   string
	GoogleMapsMapID string

	// Pembayaran (KHQR lewat ABA PayWay)
	PaymentProvider    string // payway · fake (pengembangan: tombol "simulasi bayar") · off (pra-peluncuran: checkout ditutup)
	PayWayBaseURL      string // sandbox: https://checkout-sandbox.payway.com.kh
	PayWayMerchantID   string
	PayWayAPIKey       string // rahasia: hanya di .env
	PaymentLifetimeMin int    // QR berlaku berapa menit sebelum pesanan dibatalkan

	// Kunci enkripsi data sensitif (nomor rekening seller). 32 byte, base64. Hilang = data tidak bisa dibaca.
	DataKey string
	// Cari dengan foto: alamat layanan embedder (embedder/server.py). Kosong = fitur mati.
	EmbedderURL string
}

func Load() (Config, error) {
	loadDotEnv(".env", "../.env")
	c := Config{
		Env:           get("APP_ENV", "development"),
		HTTPAddr:      get("HTTP_ADDR", ":8080"),
		PublicBaseURL: strings.TrimRight(get("PUBLIC_BASE_URL", "http://localhost:8080"), "/"),
		SellerBaseURL: strings.TrimRight(os.Getenv("SELLER_BASE_URL"), "/"),
		AdminBaseURL:  strings.TrimRight(os.Getenv("ADMIN_BASE_URL"), "/"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		DefaultLang:   get("DEFAULT_LANG", "en"),
		LogLevel:      get("LOG_LEVEL", "info"),

		JWTSecret:           os.Getenv("JWT_SECRET"),
		TelegramBotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		PushContact:         os.Getenv("PUSH_CONTACT"),
		TelegramBotUsername: strings.TrimPrefix(os.Getenv("TELEGRAM_BOT_USERNAME"), "@"),
		TelegramMode:        get("TELEGRAM_MODE", "polling"),
		TelegramWebhookKey:  os.Getenv("TELEGRAM_WEBHOOK_KEY"),
		SMSProvider:         get("SMS_PROVIDER", "log"),
		StorageDir:          get("STORAGE_DIR", "./data/uploads"),

		Translator:          get("TRANSLATOR", "claude"),
		AnthropicKey:        os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicModel:      get("ANTHROPIC_MODEL", "claude-haiku-5-5"),
		AnthropicBaseURL:    os.Getenv("ANTHROPIC_BASE_URL"), // askbot:
		GeminiKey:           os.Getenv("GEMINI_API_KEY"),     // cadangan .env; utamanya diisi di Admin → Integrasi
		GeminiModel:         get("GEMINI_MODEL", "gemini-3.5-flash"),
		GoogleTranslateKey:  os.Getenv("GOOGLE_TRANSLATE_API_KEY"),
		TranslationFallback: get("TRANSLATION_FALLBACK", "google"),
		WorkerInServe:       get("WORKER_IN_SERVE", "true") == "true",
		GoogleMapsKey:       os.Getenv("GOOGLE_MAPS_API_KEY"),
		GoogleMapsMapID:     get("GOOGLE_MAPS_MAP_ID", "DEMO_MAP_ID"),

		PaymentProvider:  get("PAYMENT_PROVIDER", "fake"),
		PayWayBaseURL:    strings.TrimRight(get("PAYWAY_BASE_URL", "https://checkout-sandbox.payway.com.kh"), "/"),
		PayWayMerchantID: os.Getenv("PAYWAY_MERCHANT_ID"),
		PayWayAPIKey:     os.Getenv("PAYWAY_API_KEY"),
	}
	c.PaymentLifetimeMin = 15
	if v, err := strconv.Atoi(os.Getenv("PAYMENT_LIFETIME_MIN")); err == nil && v >= 3 && v <= 120 {
		c.PaymentLifetimeMin = v
	}
	c.AdminMFADevSkip = c.Env == "development" && !strings.EqualFold(strings.TrimSpace(os.Getenv("ADMIN_MFA_SKIP")), "false")
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL belum diisi (lihat .env.example)")
	}
	if len(c.JWTSecret) < 32 {
		if c.Env != "development" {
			return c, fmt.Errorf("JWT_SECRET wajib diisi (minimal 32 karakter)")
		}
		c.JWTSecret = "dev-only-insecure-secret-change-me-0123456789"
	}
	c.EmbedderURL = os.Getenv("EMBEDDER_URL")
	c.DataKey = os.Getenv("DATA_ENCRYPTION_KEY")
	if c.DataKey == "" {
		if c.Env != "development" {
			return c, fmt.Errorf("DATA_ENCRYPTION_KEY wajib diisi (32 byte base64: openssl rand -base64 32)")
		}
		c.DataKey = "ZGV2LW9ubHktZGF0YS1rZXktMzItYnl0ZXMtLS0tLS0=" // hanya untuk laptop
	}
	if _, err := secretbox.New(c.DataKey); err != nil {
		return c, err
	}
	if c.TelegramBotToken == "" {
		c.TelegramMode = "off"
	} else if c.TelegramBotUsername == "" {
		return c, fmt.Errorf("TELEGRAM_BOT_USERNAME wajib diisi kalau TELEGRAM_BOT_TOKEN diisi")
	} else if c.TelegramMode == "webhook" && len(c.TelegramWebhookKey) < 16 {
		return c, fmt.Errorf("TELEGRAM_MODE=webhook butuh TELEGRAM_WEBHOOK_KEY acak minimal 16 karakter (openssl rand -hex 16)")
	}
	switch c.SMSProvider {
	case "off":
		if c.TelegramMode == "off" {
			return c, fmt.Errorf("SMS_PROVIDER=off butuh login Telegram (TELEGRAM_BOT_TOKEN) — kalau tidak, tidak ada yang bisa login")
		}
	case "log":
		if c.Env == "production" {
			return c, fmt.Errorf("SMS_PROVIDER=log tidak boleh dipakai di production (pakai off sampai ada provider SMS)")
		}
	default:
		return c, fmt.Errorf("SMS_PROVIDER harus log atau off")
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
	switch c.PaymentProvider {
	case "payway":
		if c.PayWayMerchantID == "" || c.PayWayAPIKey == "" {
			return c, fmt.Errorf("PAYMENT_PROVIDER=payway butuh PAYWAY_MERCHANT_ID dan PAYWAY_API_KEY")
		}
	case "fake":
		if c.Env == "production" {
			return c, fmt.Errorf("PAYMENT_PROVIDER=fake tidak boleh dipakai di production")
		}
	case "off":
	default:
		return c, fmt.Errorf("PAYMENT_PROVIDER harus payway, fake, atau off")
	}
	switch c.DefaultLang {
	case "en", "km", "zh":
	default:
		return c, fmt.Errorf("DEFAULT_LANG harus en, km, atau zh")
	}
	return c, nil
}

// SMSEnabled: login dengan kode SMS tersedia.
func (c Config) SMSEnabled() bool { return c.SMSProvider != "off" }

// CheckoutOpen: pembeli boleh checkout. PAYMENT_PROVIDER=off = pra-peluncuran (seller sudah bisa daftar, pembeli belum bisa beli).
func (c Config) CheckoutOpen() bool { return c.PaymentProvider != "off" }

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
