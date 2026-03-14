package domain

import "context"

type OwnerRepository interface {
	List(ctx context.Context, activeOnly bool) ([]Owner, error)
	GetByID(ctx context.Context, id int) (*Owner, error)
	GetCareerStats(ctx context.Context, ownerID int) (*OwnerCareerStats, error)
}

type SeasonRepository interface {
	GetStandings(ctx context.Context, year int) ([]SeasonStanding, error)
	GetAwards(ctx context.Context, year int) ([]SeasonAwardWithDetails, error)
}

type BeltRepository interface {
	GetCurrent(ctx context.Context) (*BeltHistoryWithDetails, error)
	GetLeaderboard(ctx context.Context) ([]OwnerBeltWeeks, error)
}
