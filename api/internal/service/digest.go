package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
	"github.com/rknruben56/juwanna-fantasy/api/internal/sleeper"
)

// DigestService combines data from multiple services into a weekly digest.
type DigestService struct {
	sleeperSvc   *SleeperService
	analyticsSvc *AnalyticsService
	rivalries    domain.RivalryRepository
	client       *sleeper.Client
}

// NewDigestService creates a new DigestService.
func NewDigestService(
	sleeperSvc *SleeperService,
	analyticsSvc *AnalyticsService,
	rivalries domain.RivalryRepository,
	client *sleeper.Client,
) *DigestService {
	return &DigestService{
		sleeperSvc:   sleeperSvc,
		analyticsSvc: analyticsSvc,
		rivalries:    rivalries,
		client:       client,
	}
}

// GetWeeklyDigest returns a combined weekly summary with matchups, standings,
// belt status, award projections, and power rankings.
func (s *DigestService) GetWeeklyDigest(ctx context.Context) (*domain.WeeklyDigest, error) {
	// Get NFL state for season/week metadata
	leagueID, err := s.sleeperSvc.GetLeagueID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get league id: %w", err)
	}

	state, err := s.client.GetNFLState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get nfl state: %w", err)
	}

	digest := &domain.WeeklyDigest{
		Season: state.Season,
		Week:   state.Week,
	}

	// Fetch data in parallel
	var (
		matchups   []domain.LiveMatchup
		standings  []domain.LiveStanding
		belt       *domain.LiveBeltStatus
		awards     []domain.AwardProjection
		rankings   []domain.PowerRanking
	)

	g, gctx := errgroup.WithContext(ctx)
	_ = leagueID // used indirectly by sleeper service

	// Required: matchups
	g.Go(func() error {
		var err error
		matchups, err = s.sleeperSvc.GetLiveMatchups(gctx)
		if err != nil {
			return fmt.Errorf("get live matchups: %w", err)
		}
		return nil
	})

	// Required: standings
	g.Go(func() error {
		var err error
		standings, err = s.sleeperSvc.GetLiveStandings(gctx)
		if err != nil {
			return fmt.Errorf("get live standings: %w", err)
		}
		return nil
	})

	// Optional: belt status
	g.Go(func() error {
		var err error
		belt, err = s.sleeperSvc.GetLiveBeltStatus(gctx)
		if err != nil {
			log.Warn().Err(err).Msg("digest: failed to get belt status")
		}
		return nil
	})

	// Optional: award projections
	g.Go(func() error {
		var err error
		awards, err = s.analyticsSvc.GetAwardProjections(gctx)
		if err != nil {
			log.Warn().Err(err).Msg("digest: failed to get award projections")
		}
		return nil
	})

	// Optional: power rankings
	g.Go(func() error {
		var err error
		rankings, err = s.analyticsSvc.GetPowerRankings(gctx)
		if err != nil {
			log.Warn().Err(err).Msg("digest: failed to get power rankings")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Enrich matchups with H2H context
	digest.Matchups = s.enrichMatchups(ctx, matchups)
	digest.Standings = standings
	digest.BeltStatus = belt
	digest.AwardProjections = awards
	digest.PowerRankings = rankings

	return digest, nil
}

// enrichMatchups adds head-to-head context to matchups where both teams have mapped owners.
func (s *DigestService) enrichMatchups(ctx context.Context, matchups []domain.LiveMatchup) []domain.DigestMatchup {
	result := make([]domain.DigestMatchup, len(matchups))

	for i, m := range matchups {
		result[i] = domain.DigestMatchup{LiveMatchup: m}

		if len(m.Teams) != 2 {
			continue
		}

		t1 := m.Teams[0]
		t2 := m.Teams[1]

		if t1.OwnerID == nil || t2.OwnerID == nil {
			continue
		}

		h2h, err := s.rivalries.GetHeadToHead(ctx, *t1.OwnerID, *t2.OwnerID)
		if err != nil {
			log.Debug().Err(err).
				Int("owner1", *t1.OwnerID).
				Int("owner2", *t2.OwnerID).
				Msg("digest: no h2h data")
			continue
		}

		result[i].HeadToHead = &domain.DigestH2H{
			OwnerWins:    h2h.Wins,
			OpponentWins: h2h.Losses,
			TotalGames:   h2h.TotalGames,
		}
	}

	return result
}
