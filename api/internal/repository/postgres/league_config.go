package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type LeagueConfigRepo struct {
	pool *pgxpool.Pool
}

func NewLeagueConfigRepo(pool *pgxpool.Pool) *LeagueConfigRepo {
	return &LeagueConfigRepo{pool: pool}
}

func (r *LeagueConfigRepo) GetBySeasonID(ctx context.Context, seasonID int) (*domain.SleeperLeagueConfig, error) {
	var c domain.SleeperLeagueConfig
	err := r.pool.QueryRow(ctx, `
		SELECT id, season_id, sleeper_league_id, created_at
		FROM sleeper_league_config
		WHERE season_id = $1
	`, seasonID).Scan(&c.ID, &c.SeasonID, &c.SleeperLeagueID, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get league config: %w", err)
	}
	return &c, nil
}

func (r *LeagueConfigRepo) GetCurrentLeagueID(ctx context.Context) (string, error) {
	var leagueID string
	err := r.pool.QueryRow(ctx, `
		SELECT lc.sleeper_league_id
		FROM sleeper_league_config lc
		JOIN seasons s ON lc.season_id = s.id
		ORDER BY s.year DESC
		LIMIT 1
	`).Scan(&leagueID)
	if err != nil {
		return "", fmt.Errorf("get current league id: %w", err)
	}
	return leagueID, nil
}

func (r *LeagueConfigRepo) Set(ctx context.Context, seasonID int, sleeperLeagueID string) (*domain.SleeperLeagueConfig, error) {
	var c domain.SleeperLeagueConfig
	err := r.pool.QueryRow(ctx, `
		INSERT INTO sleeper_league_config (season_id, sleeper_league_id)
		VALUES ($1, $2)
		ON CONFLICT (season_id)
		DO UPDATE SET sleeper_league_id = EXCLUDED.sleeper_league_id
		RETURNING id, season_id, sleeper_league_id, created_at
	`, seasonID, sleeperLeagueID).Scan(&c.ID, &c.SeasonID, &c.SleeperLeagueID, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("set league config: %w", err)
	}
	return &c, nil
}
