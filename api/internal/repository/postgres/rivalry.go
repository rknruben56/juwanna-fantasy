package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type RivalryRepo struct {
	pool *pgxpool.Pool
}

func NewRivalryRepo(pool *pgxpool.Pool) *RivalryRepo {
	return &RivalryRepo{pool: pool}
}

func (r *RivalryRepo) GetHeadToHead(ctx context.Context, ownerID, opponentID int) (*domain.HeadToHeadRecord, error) {
	var h domain.HeadToHeadRecord
	err := r.pool.QueryRow(ctx,
		`SELECT owner_id, owner_name, opponent_id, opponent_name,
			wins, losses, total_games, win_pct,
			total_points_for, total_points_against
		FROM rivalry_stats
		WHERE owner_id = $1 AND opponent_id = $2`,
		ownerID, opponentID,
	).Scan(
		&h.OwnerID, &h.OwnerName, &h.OpponentID, &h.OpponentName,
		&h.Wins, &h.Losses, &h.TotalGames, &h.WinPct,
		&h.TotalPointsFor, &h.TotalPointsAgainst,
	)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *RivalryRepo) GetAllRivalries(ctx context.Context, ownerID int) ([]domain.HeadToHeadRecord, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT owner_id, owner_name, opponent_id, opponent_name,
			wins, losses, total_games, win_pct,
			total_points_for, total_points_against
		FROM rivalry_stats
		WHERE owner_id = $1
		ORDER BY total_games DESC, win_pct DESC`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.HeadToHeadRecord
	for rows.Next() {
		var h domain.HeadToHeadRecord
		if err := rows.Scan(
			&h.OwnerID, &h.OwnerName, &h.OpponentID, &h.OpponentName,
			&h.Wins, &h.Losses, &h.TotalGames, &h.WinPct,
			&h.TotalPointsFor, &h.TotalPointsAgainst,
		); err != nil {
			return nil, err
		}
		records = append(records, h)
	}
	return records, rows.Err()
}

func (r *RivalryRepo) GetMatchupHistory(ctx context.Context, ownerID, opponentID int) ([]domain.MatchupHistoryEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.year, mh.week_number, o1.name, o2.name,
			mh.points_scored, mh.opponent_points, mh.is_playoff
		FROM matchup_history mh
		JOIN seasons s ON mh.season_id = s.id
		JOIN owners o1 ON mh.owner_id = o1.id
		JOIN owners o2 ON mh.opponent_id = o2.id
		WHERE mh.owner_id = $1 AND mh.opponent_id = $2
		ORDER BY s.year DESC, mh.week_number DESC`,
		ownerID, opponentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.MatchupHistoryEntry
	for rows.Next() {
		var e domain.MatchupHistoryEntry
		if err := rows.Scan(
			&e.Season, &e.Week, &e.OwnerName, &e.OpponentName,
			&e.PointsScored, &e.OpponentPoints, &e.IsPlayoff,
		); err != nil {
			return nil, err
		}
		if e.PointsScored > e.OpponentPoints {
			e.Result = "W"
		} else if e.PointsScored < e.OpponentPoints {
			e.Result = "L"
		} else {
			e.Result = "T"
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
