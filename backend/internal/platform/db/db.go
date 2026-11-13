// Package db membuka koneksi pool PostgreSQL (pgx).
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL tidak valid: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	cfg.ConnConfig.RuntimeParams["application_name"] = "phsar"

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("tidak bisa terhubung ke database: %w", err)
	}
	return pool, nil
}

// WithActor menjalankan fn di dalam transaksi dan mengisi app.actor_id / app.actor_type /
// app.request_id, yang dibaca trigger audit_logs. Semua perubahan data dari handler
// sebaiknya lewat fungsi ini supaya tercatat siapa yang mengubah.
func WithActor(ctx context.Context, pool *pgxpool.Pool, actorID, actorType, requestID string, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx,
		`SELECT set_config('app.actor_id', $1, true),
		        set_config('app.actor_type', $2, true),
		        set_config('app.request_id', $3, true)`,
		actorID, actorType, requestID); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
