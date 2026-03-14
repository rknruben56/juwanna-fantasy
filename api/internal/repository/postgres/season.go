package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type SeasonRepo struct {
	pool *pgxpool.Pool
}

func NewSeasonRepo(pool *pgxpool.Pool) *SeasonRepo {
	return &SeasonRepo{pool: pool}
}

func (r *SeasonRepo) GetStandings(ctx context.Context, year int) ([]domain.SeasonStanding, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT year, owner_name, regular_season_wins, regular_season_losses,
			points_for, points_against, made_playoffs, team_name
		FROM season_leaderboard WHERE year = $1`, year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var standings []domain.SeasonStanding
	for rows.Next() {
		var s domain.SeasonStanding
		if err := rows.Scan(&s.Year, &s.OwnerName, &s.RegularSeasonWins, &s.RegularSeasonLosses,
			&s.PointsFor, &s.PointsAgainst, &s.MadePlayoffs, &s.TeamName); err != nil {
			return nil, err
		}
		standings = append(standings, s)
	}
	return standings, rows.Err()
}

func (r *SeasonRepo) GetAwards(ctx context.Context, year int) ([]domain.SeasonAwardWithDetails, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT sa.id, o.name, at.name, s.year, sa.stat_value, sa.team_name
		FROM season_awards sa
		JOIN owners o ON sa.owner_id = o.id
		JOIN seasons s ON sa.season_id = s.id
		JOIN award_types at ON sa.award_type_id = at.id
		WHERE s.year = $1
		ORDER BY at.name`, year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []domain.SeasonAwardWithDetails
	for rows.Next() {
		var a domain.SeasonAwardWithDetails
		if err := rows.Scan(&a.ID, &a.OwnerName, &a.AwardName, &a.Year, &a.StatValue, &a.TeamName); err != nil {
			return nil, err
		}
		awards = append(awards, a)
	}
	return awards, rows.Err()
}
