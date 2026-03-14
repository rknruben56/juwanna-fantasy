package domain

import (
	"context"
	"time"
)

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

type MappingRepository interface {
	List(ctx context.Context) ([]OwnerSleeperMappingWithName, error)
	GetByOwnerID(ctx context.Context, ownerID int) (*OwnerSleeperMapping, error)
	GetBySleeperUserID(ctx context.Context, sleeperUserID string) (*OwnerSleeperMapping, error)
	Create(ctx context.Context, ownerID int, sleeperUserID string) (*OwnerSleeperMapping, error)
	Update(ctx context.Context, id int, sleeperUserID string) (*OwnerSleeperMapping, error)
	Delete(ctx context.Context, id int) error
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, expiresAt time.Time) error
	Delete(ctx context.Context, key string) error
	DeleteExpired(ctx context.Context) error
}

type LeagueConfigRepository interface {
	GetBySeasonID(ctx context.Context, seasonID int) (*SleeperLeagueConfig, error)
	GetCurrentLeagueID(ctx context.Context) (string, error)
	Set(ctx context.Context, seasonID int, sleeperLeagueID string) (*SleeperLeagueConfig, error)
}

type RivalryRepository interface {
	GetHeadToHead(ctx context.Context, ownerID, opponentID int) (*HeadToHeadRecord, error)
	GetAllRivalries(ctx context.Context, ownerID int) ([]HeadToHeadRecord, error)
	GetMatchupHistory(ctx context.Context, ownerID, opponentID int) ([]MatchupHistoryEntry, error)
}
