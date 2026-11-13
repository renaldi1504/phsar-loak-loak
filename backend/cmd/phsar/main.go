// Command phsar adalah satu binary untuk seluruh backend Phsar.
//
//	phsar serve            jalankan HTTP API (/v1/app, /v1/seller, /admin, /webhooks)
//	phsar migrate up       jalankan migrasi database
//	phsar migrate down     mundurkan satu migrasi
//	phsar migrate status   lihat status migrasi
//	phsar worker           job queue (River) — belum diaktifkan di Tahap 0
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/config"
	"phsar/internal/httpx"
	"phsar/internal/platform/db"
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

	case "worker":
		log.Info("worker belum diaktifkan; akan memakai River di Tahap 1 (terjemahan AI, transcode video)")
		<-ctx.Done()
		return nil

	default:
		return usage()
	}
}

func usage() error {
	return fmt.Errorf("pemakaian: phsar <serve|migrate [up|down|status|reset]|admin grant <hp> <peran>|worker>")
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
