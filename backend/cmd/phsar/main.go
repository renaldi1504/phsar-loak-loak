// Command phsar adalah satu binary untuk seluruh backend Phsar.
//
//	phsar serve            jalankan HTTP API (/v1/app, /v1/seller, /admin, /webhooks)
//	phsar migrate up       jalankan migrasi database
//	phsar migrate down     mundurkan satu migrasi
//	phsar migrate status   lihat status migrasi
//	phsar worker           worker terjemahan AI (dev: sudah ikut jalan di serve)
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/config"
	"phsar/internal/httpx"
	"phsar/internal/media"
	"phsar/internal/payment"
	"phsar/internal/platform/db"
	"phsar/internal/translate"
	"phsar/migrations"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := newLogger(cfg)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch args[0] {
	case "serve":
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		return httpx.Serve(ctx, cfg, log, pool)

	case "migrate":
		sub := "up"
		if len(args) > 1 {
			sub = args[1]
		}
		return migrations.Run(ctx, cfg.DatabaseURL, sub)

	case "admin":
		// phsar admin grant <nomor HP | user_id> <peran>   ·   phsar admin reset-mfa <nomor HP | user_id>
		usage := fmt.Errorf("pemakaian: phsar admin grant <nomor HP | @username Telegram | user_id> <super_admin|moderator|cs|finance|ads_reviewer>\n" +
			"         phsar admin reset-mfa <nomor HP | user_id>")
		if len(args) < 3 {
			return usage
		}
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		switch {
		case args[1] == "grant" && len(args) == 4:
			return grantAdmin(ctx, pool, args[2], args[3])
		case args[1] == "reset-mfa" && len(args) == 3:
			return resetAdminMFA(ctx, pool, args[2])
		}
		return usage

	case "translate-test":
		// phsar translate-test ["teks"] — cek API key & model terjemahan tanpa membuka Seller Center.
		text := "Cotton T-shirt, red, size M. Soft and comfortable."
		if len(args) > 1 {
			text = strings.Join(args[1:], " ")
		}
		return translateTest(ctx, cfg, text)

	case "media-backfill":
		// Antrekan video/foto lama (sebelum ada kompres) supaya ikut diproses worker.
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		n, err := media.Backfill(ctx, pool)
		if err != nil {
			return err
		}
		fmt.Printf("%d media diantrekan. Jalankan `make run` (atau `phsar worker`) untuk memprosesnya.\n", n)
		return nil

	case "seed-demo":
		// phsar seed-demo [--reset] — toko & produk contoh untuk mencoba website (jangan di production).
		if cfg.Env == "production" && !(len(args) > 1 && args[len(args)-1] == "--force") {
			return fmt.Errorf("seed-demo menolak berjalan di production (tambahkan --force kalau benar-benar perlu)")
		}
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		sd := httpx.NewDemoSeeder(cfg, log, pool)
		if len(args) > 1 && args[1] == "--reset" {
			n, err := sd.Reset(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("%d toko demo (dan produknya) disembunyikan.\n", n)
			return nil
		}
		res, err := sd.Seed(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("✅ %d toko & %d produk demo dibuat (%d toko sudah ada, dilewati).\n", res.Stores, res.Products, res.Skipped)
		if !res.Videos {
			fmt.Println("ℹ️  ffmpeg tidak ditemukan: produk demo dibuat tanpa video (pasang: brew install ffmpeg, lalu --reset dan jalankan lagi).")
		}
		fmt.Println("Buka http://localhost:8080/ untuk melihatnya. Hapus lagi dengan: phsar seed-demo --reset")
		return nil

	case "payway-test":
		// phsar payway-test            buat QR $0.01 di PayWay (kunci dari .env) → simpan gambar QR, lalu cek statusnya
		// phsar payway-test <tran_id>  cek status transaksi uji yang sudah dibayar
		return paywayTest(ctx, cfg, args[1:])

	case "worker":
		// Worker terpisah untuk production (set WORKER_IN_SERVE=false di server API).
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		httpx.StartWorkers(ctx, cfg, log, pool)
		<-ctx.Done()
		return nil

	default:
		return usage()
	}
}

func translateTest(ctx context.Context, cfg config.Config, text string) error {
	t := httpx.NewTranslator(cfg)
	if t == nil {
		return fmt.Errorf("terjemahan nonaktif: ANTHROPIC_API_KEY kosong (TRANSLATOR=%s). Isi di .env lalu jalankan lagi", cfg.Translator)
	}
	keyHint := "(kosong)"
	if k := cfg.AnthropicKey; len(k) > 12 {
		keyHint = k[:10] + "…" + k[len(k)-4:]
	}
	fmt.Printf("Provider: %s · model: %s · key: %s\nTeks: %s\n\n", cfg.Translator, cfg.AnthropicModel, keyHint, text)
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	res, err := t.Translate(ctx, translate.Request{Kind: "product", Source: "en", Targets: []string{"km", "zh"}, Fields: map[string]string{"name": text}})
	if err != nil {
		return fmt.Errorf("GAGAL: %w", err)
	}
	fmt.Printf("✅ Berhasil (model %s)\n  km: %s\n  zh: %s\n", res.Model, res.Texts["km"]["name"], res.Texts["zh"]["name"])
	return nil
}

func usage() error {
	return fmt.Errorf("pemakaian: phsar <serve|migrate [up|down|status|reset]|admin grant <hp> <peran>|admin reset-mfa <hp>|worker|translate-test|payway-test [tran_id]|media-backfill|seed-demo [--reset]>")
}

// paywayTest: uji kunci PayWay tanpa database & tanpa membuka website. Kunci tidak pernah dicetak.
func paywayTest(ctx context.Context, cfg config.Config, args []string) error {
	if cfg.PayWayMerchantID == "" || cfg.PayWayAPIKey == "" {
		return fmt.Errorf("isi PAYWAY_MERCHANT_ID dan PAYWAY_API_KEY di .env dulu (PAYWAY_BASE_URL sekarang: %s)", cfg.PayWayBaseURL)
	}
	p := &payment.PayWay{BaseURL: cfg.PayWayBaseURL, MerchantID: cfg.PayWayMerchantID, APIKey: cfg.PayWayAPIKey}
	fmt.Printf("PayWay: %s · merchant %s\n", cfg.PayWayBaseURL, cfg.PayWayMerchantID)
	tranID := ""
	if len(args) > 0 {
		tranID = args[0]
	} else {
		tranID = fmt.Sprintf("T%d", time.Now().Unix())
		qr, err := p.CreateQR(ctx, payment.QRRequest{TranID: tranID, Amount: "0.01", FirstName: "Phsar test", LifetimeMin: 10,
			Items: []payment.Item{{Name: "Phsar test", Quantity: 1, Price: "0.01"}}, CallbackURL: cfg.PublicBaseURL + "/webhooks/payway"})
		if err != nil {
			return fmt.Errorf("%w\n  kode 1 = API key salah · 6/102 = domain/IP belum di-whitelist ABA · 23/32 = QR API belum diaktifkan untuk merchant ini (minta ke ABA)", err)
		}
		fmt.Printf("✅ QR dibuat untuk $0.01 — tran_id %s\n", tranID)
		if img, ok := strings.CutPrefix(qr.Image, "data:image/png;base64,"); ok {
			if b, err := base64.StdEncoding.DecodeString(img); err == nil && os.WriteFile("payway-test-qr.png", b, 0o644) == nil {
				fmt.Println("   Gambar QR: payway-test-qr.png — buka lalu scan dengan ABA Mobile / aplikasi bank (sandbox: pakai akun uji ABA).")
			}
		}
		if qr.Deeplink != "" {
			fmt.Println("   Link ABA Mobile:", qr.Deeplink)
		}
	}
	st, err := p.Check(ctx, tranID)
	if err != nil {
		return err
	}
	fmt.Printf("Status %s: %s", tranID, st.State)
	if st.Amount != "" && st.Amount != "<nil>" {
		fmt.Printf(" · $%s", st.Amount)
	}
	fmt.Println()
	if !st.Paid() {
		fmt.Printf("   Setelah scan & bayar, cek lagi:  phsar payway-test %s\n", tranID)
	}
	return nil
}

func newLogger(cfg config.Config) *slog.Logger {
	var level slog.Level
	_ = level.UnmarshalText([]byte(cfg.LogLevel))
	opts := &slog.HandlerOptions{Level: level}
	if cfg.Env == "development" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func findUser(ctx context.Context, pool *pgxpool.Pool, who string) (string, string, error) {
	userID, who, err := auth.FindUser(ctx, pool, who)
	if err != nil {
		return "", who, fmt.Errorf("user %s tidak ditemukan (pakai nomor HP, @username Telegram, atau tg:<id>; orangnya harus login sekali di /admin dulu): %w", who, err)
	}
	return userID, who, nil
}

// resetAdminMFA: admin kehilangan HP → MFA dihapus, dipasang lagi saat login berikutnya.
func resetAdminMFA(ctx context.Context, pool *pgxpool.Pool, who string) error {
	userID, who, err := findUser(ctx, pool, who)
	if err != nil {
		return err
	}
	return db.WithActor(ctx, pool, "", "system", "cli", func(tx pgx.Tx) error {
		if err := auth.ResetMFA(ctx, tx, "", userID); err != nil {
			return fmt.Errorf("%s bukan admin: %w", who, err)
		}
		fmt.Printf("OK: MFA %s dihapus. Saat login berikutnya di /admin, ia memasang aplikasi authenticator lagi.\n", who)
		return nil
	})
}

// grantAdmin menjadikan user (yang sudah pernah login) sebagai admin Phsar.
// Semua admin wajib memasang MFA (aplikasi authenticator) saat pertama login di /admin.
func grantAdmin(ctx context.Context, pool *pgxpool.Pool, who, role string) error {
	userID, who, err := findUser(ctx, pool, who)
	if err != nil {
		return err
	}
	return db.WithActor(ctx, pool, "", "system", "cli", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO admin_users (user_id, role_code) VALUES ($1, $2)
			ON CONFLICT (user_id) DO UPDATE SET role_code = EXCLUDED.role_code, active = true, deleted_at = NULL`, userID, role)
		if err == nil {
			fmt.Printf("OK: %s sekarang admin dengan peran %s. Login di /admin lalu pasang aplikasi authenticator (MFA).\n", who, role)
		}
		return err
	})
}
