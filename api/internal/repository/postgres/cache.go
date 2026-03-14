package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CacheRepo struct {
	pool *pgxpool.Pool
}

func NewCacheRepo(pool *pgxpool.Pool) *CacheRepo {
	return &CacheRepo{pool: pool}
}

func (r *CacheRepo) Get(ctx context.Context, key string) ([]byte, error) {
	var data []byte
	err := r.pool.QueryRow(ctx, `
		SELECT data FROM sleeper_cache
		WHERE cache_key = $1 AND expires_at > NOW()
	`, key).Scan(&data)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache get: %w", err)
	}
	return data, nil
}

func (r *CacheRepo) Set(ctx context.Context, key string, data []byte, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sleeper_cache (cache_key, data, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (cache_key)
		DO UPDATE SET data = EXCLUDED.data, expires_at = EXCLUDED.expires_at, updated_at = CURRENT_TIMESTAMP
	`, key, data, expiresAt)
	if err != nil {
		return fmt.Errorf("cache set: %w", err)
	}
	return nil
}

func (r *CacheRepo) Delete(ctx context.Context, key string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sleeper_cache WHERE cache_key = $1`, key)
	if err != nil {
		return fmt.Errorf("cache delete: %w", err)
	}
	return nil
}

func (r *CacheRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sleeper_cache WHERE expires_at <= NOW()`)
	if err != nil {
		return fmt.Errorf("cache delete expired: %w", err)
	}
	return nil
}
