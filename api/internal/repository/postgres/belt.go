package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type BeltRepo struct {
	pool *pgxpool.Pool
}

func NewBeltRepo(pool *pgxpool.Pool) *BeltRepo {
	return &BeltRepo{pool: pool}
}

func (r *BeltRepo) GetCurrent(ctx context.Context) (*domain.BeltHistoryWithDetails, error) {
	var b domain.BeltHistoryWithDetails
	err := r.pool.QueryRow(ctx,
		`SELECT o.name, s.year, cbh.week_number
		FROM championship_belt_history cbh
		JOIN owners o ON cbh.owner_id = o.id
		JOIN seasons s ON cbh.season_id = s.id
		ORDER BY s.year DESC, cbh.week_number DESC
		LIMIT 1`,
	).Scan(&b.OwnerName, &b.Year, &b.WeekNumber)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BeltRepo) GetLeaderboard(ctx context.Context) ([]domain.OwnerBeltWeeks, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT owner_id, name, weeks_with_belt
		FROM owner_belt_weeks
		WHERE weeks_with_belt > 0
		ORDER BY weeks_with_belt DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []domain.OwnerBeltWeeks
	for rows.Next() {
		var b domain.OwnerBeltWeeks
		if err := rows.Scan(&b.OwnerID, &b.Name, &b.WeeksWithBelt); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, b)
	}
	return leaderboard, rows.Err()
}
