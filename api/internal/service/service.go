package service

import (
	"context"

	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type HistoricalService struct {
	owners  domain.OwnerRepository
	seasons domain.SeasonRepository
	belt    domain.BeltRepository
}

func NewHistoricalService(owners domain.OwnerRepository, seasons domain.SeasonRepository, belt domain.BeltRepository) *HistoricalService {
	return &HistoricalService{owners: owners, seasons: seasons, belt: belt}
}

func (s *HistoricalService) ListOwners(ctx context.Context, activeOnly bool) ([]domain.Owner, error) {
	return s.owners.List(ctx, activeOnly)
}

func (s *HistoricalService) GetOwnerStats(ctx context.Context, ownerID int) (*domain.OwnerCareerStats, error) {
	return s.owners.GetCareerStats(ctx, ownerID)
}

func (s *HistoricalService) GetSeasonStandings(ctx context.Context, year int) ([]domain.SeasonStanding, error) {
	return s.seasons.GetStandings(ctx, year)
}

func (s *HistoricalService) GetSeasonAwards(ctx context.Context, year int) ([]domain.SeasonAwardWithDetails, error) {
	return s.seasons.GetAwards(ctx, year)
}

func (s *HistoricalService) GetCurrentBeltHolder(ctx context.Context) (*domain.BeltHistoryWithDetails, error) {
	return s.belt.GetCurrent(ctx)
}

func (s *HistoricalService) GetBeltLeaderboard(ctx context.Context) ([]domain.OwnerBeltWeeks, error) {
	return s.belt.GetLeaderboard(ctx)
}
