package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type OwnerRepo struct {
	pool *pgxpool.Pool
}

func NewOwnerRepo(pool *pgxpool.Pool) *OwnerRepo {
	return &OwnerRepo{pool: pool}
}

func (r *OwnerRepo) List(ctx context.Context, activeOnly bool) ([]domain.Owner, error) {
	query := `SELECT id, name, join_year, leave_year, created_at, updated_at FROM owners`
	if activeOnly {
		query += ` WHERE leave_year IS NULL`
	}
	query += ` ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []domain.Owner
	for rows.Next() {
		var o domain.Owner
		if err := rows.Scan(&o.ID, &o.Name, &o.JoinYear, &o.LeaveYear, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, rows.Err()
}

func (r *OwnerRepo) GetByID(ctx context.Context, id int) (*domain.Owner, error) {
	var o domain.Owner
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, join_year, leave_year, created_at, updated_at FROM owners WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.JoinYear, &o.LeaveYear, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OwnerRepo) GetCareerStats(ctx context.Context, ownerID int) (*domain.OwnerCareerStats, error) {
	var s domain.OwnerCareerStats
	err := r.pool.QueryRow(ctx,
		`SELECT owner_id, name, join_year, leave_year, years_in_league,
			total_regular_season_wins, total_regular_season_losses,
			total_playoff_wins, total_playoff_losses,
			playoff_appearances, championships, first_place_loser_trophies
		FROM owner_career_stats WHERE owner_id = $1`, ownerID,
	).Scan(
		&s.OwnerID, &s.Name, &s.JoinYear, &s.LeaveYear, &s.YearsInLeague,
		&s.TotalRegularSeasonWins, &s.TotalRegularSeasonLosses,
		&s.TotalPlayoffWins, &s.TotalPlayoffLosses,
		&s.PlayoffAppearances, &s.Championships, &s.FirstPlaceLoserTrophies,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
