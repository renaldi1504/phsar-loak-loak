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
		// phsar admin grant <nomor HP | user_id> <peran>
		if len(args) != 4 || args[1] != "grant" {
			return fmt.Errorf("pemakaian: phsar admin grant <nomor HP | user_id> <super_admin|moderator|cs|finance|ads_reviewer>")
		}
		pool, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		return grantAdmin(ctx, pool, args[2], args[3])

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
	return fmt.Errorf("pemakaian: phsar <serve|migrate [up|down|status|reset]|admin grant <hp> <peran>|worker|translate-test|media-backfill>")
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

// grantAdmin menjadikan user (yang sudah pernah login) sebagai admin Phsar.
// Peran finance & super_admin wajib MFA; sebelum login admin + MFA selesai, pakai moderator/cs.
func grantAdmin(ctx context.Context, pool *pgxpool.Pool, who, role string) error {
	q := `SELECT id FROM users WHERE id::text = $1 AND deleted_at IS NULL`
	if phone, err := auth.NormalizePhone(who); err == nil {
		who, q = phone, `SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL`
	}
	var userID string
	if err := pool.QueryRow(ctx, q, who).Scan(&userID); err != nil {
		return fmt.Errorf("user %s tidak ditemukan (login dulu sekali lewat app / /dev/login): %w", who, err)
	}
	return db.WithActor(ctx, pool, "", "system", "cli", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO admin_users (user_id, role_code) VALUES ($1, $2)
			ON CONFLICT (user_id) DO UPDATE SET role_code = EXCLUDED.role_code, active = true, deleted_at = NULL`, userID, role)
		if err == nil {
			fmt.Printf("OK: %s sekarang admin dengan peran %s\n", who, role)
		}
		return err
	})
}
