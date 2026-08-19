package database

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

func (p *Postgres) Migrate(ctx context.Context) error {
	db := stdlib.OpenDBFromPool(p.psql)
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close connection to db after migratio", "err", err)
		}
	}()

	goose.SetBaseFS(migrationFiles)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set db dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "sql"); err != nil {
		return fmt.Errorf("failed to goose up: %w", err)
	}

	return nil
}
