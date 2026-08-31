package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectionTimeout time.Duration
}

type Postgres struct {
	psql *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connString string, conf Config) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConns = conf.MaxConns
	config.MinConns = conf.MinConns
	config.MaxConnLifetime = conf.MaxConnLifeTime
	config.MaxConnIdleTime = conf.MaxConnIdleTime
	config.HealthCheckPeriod = conf.HealthCheckPeriod
	config.ConnConfig.ConnectTimeout = conf.ConnectionTimeout

	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": "url-shortener",
		"timezone":         "UTC",
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer()

	db, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{
		psql: db,
	}, nil
}

func (p *Postgres) Set(ctx context.Context, url *domain.URL) error {
	query := `
			INSERT INTO urls (id, user_id, long_url, expired_at)
			VALUES ($1, $2, $3, $4)
	`
	_, err := p.psql.Exec(ctx, query, url.ID, url.UserID, url.LongURL, url.ExpiredAt)

	if err != nil {
		return fmt.Errorf("failed to insert url: %w", err)
	}

	return nil
}

func (p *Postgres) Get(ctx context.Context, id uint64) (domain.URL, error) {
	query := `
		SELECT id, user_id, long_url, expired_at
		FROM urls
		WHERE id = $1 AND expired_at > NOW()
	`

	var url domain.URL
	err := p.psql.QueryRow(ctx, query, id).Scan(
		&url.ID,
		&url.UserID,
		&url.LongURL,
		&url.ExpiredAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.URL{}, domain.ErrUrlNotFound
		}
		return domain.URL{}, fmt.Errorf("failed to get url for %d: %w", id, err)
	}

	return url, nil
}

func (p *Postgres) Close() {
	p.psql.Close()
}
