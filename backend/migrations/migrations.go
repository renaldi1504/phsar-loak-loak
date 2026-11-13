// Package migrations menyimpan file SQL (goose) di dalam binary.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

func Run(ctx context.Context, databaseURL, cmd string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	goose.SetBaseFS(files)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	switch cmd {
	case "up":
		return goose.UpContext(ctx, db, ".")
	case "down":
		return goose.DownContext(ctx, db, ".")
	case "status":
		return goose.StatusContext(ctx, db, ".")
	case "reset":
		return goose.ResetContext(ctx, db, ".")
	default:
		return fmt.Errorf("perintah migrate tidak dikenal: %s", cmd)
	}
}
