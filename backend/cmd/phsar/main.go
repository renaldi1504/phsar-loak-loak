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

	case "worker":
		log.Info("worker belum diaktifkan; akan memakai River di Tahap 1 (terjemahan AI, transcode video)")
		<-ctx.Done()
		return nil

	default:
		return usage()
	}
}

func usage() error {
	return fmt.Errorf("pemakaian: phsar <serve|migrate [up|down|status|reset]|worker>")
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
