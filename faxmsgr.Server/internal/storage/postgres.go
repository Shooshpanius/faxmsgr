// Пакет storage предоставляет клиентов для PostgreSQL и Redis.
package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres создаёт пул соединений к PostgreSQL по DSN-строке.
func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
