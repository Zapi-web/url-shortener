package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

func (p *Postgres) Migrate(ctx context.Context) error {
	db := stdlib.OpenDBFromPool(p.psql)
	defer db.Close()

	goose.SetBaseFS(migrationFiles)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set db dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "sql"); err != nil {
		return fmt.Errorf("failed to goose up: %w", err)
	}

	return nil
}
