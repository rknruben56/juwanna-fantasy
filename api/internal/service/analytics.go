package service

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
	"github.com/rknruben56/juwanna-fantasy/api/internal/sleeper"
)

// AnalyticsService handles rivalries, power rankings, and award projections.
type AnalyticsService struct {
	rivalries       domain.RivalryRepository
	sleeperClient   *sleeper.Client
	sleeperSvc      *SleeperService
}

// NewAnalyticsService creates a new AnalyticsService.
func NewAnalyticsService(
	rivalries domain.RivalryRepository,
	sleeperClient *sleeper.Client,
	sleeperSvc *SleeperService,
) *AnalyticsService {
	return &AnalyticsService{
		rivalries:     rivalries,
		sleeperClient: sleeperClient,
		sleeperSvc:    sleeperSvc,
	}
}

// GetHeadToHead returns the head-to-head record between two owners.
func (s *AnalyticsService) GetHeadToHead(ctx context.Context, ownerID, opponentID int) (*domain.RivalrySummary, error) {
	record, err := s.rivalries.GetHeadToHead(ctx, ownerID, opponentID)
	if err != nil {
		return nil, fmt.Errorf("get h2h record: %w", err)
	}

	matchups, err := s.rivalries.GetMatchupHistory(ctx, ownerID, opponentID)
	if err != nil {
		return nil, fmt.Errorf("get matchup history: %w", err)
	}

	// Limit to last 10 recent matchups
	limit := 10
	if len(matchups) < limit {
		limit = len(matchups)
	}

	return &domain.RivalrySummary{
		HeadToHeadRecord: *record,
		RecentMatchups:   matchups[:limit],
	}, nil
}

// GetAllRivalries returns all head-to-head records for a given owner.
func (s *AnalyticsService) GetAllRivalries(ctx context.Context, ownerID int) ([]domain.HeadToHeadRecord, error) {
	return s.rivalries.GetAllRivalries(ctx, ownerID)
}

// GetPowerRankings computes power rankings from live Sleeper data.
// Score = (WinPct * 60) + (NormalizedPointsFor * 30) + (StreakBonus * 10)
func (s *AnalyticsService) GetPowerRankings(ctx context.Context) ([]domain.PowerRanking, error) {
	standings, err := s.sleeperSvc.GetLiveStandings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get live standings: %w", err)
	}

	if len(standings) == 0 {
		return nil, nil
	}

	// Find max points for normalization
	var maxPF float64
	for _, st := range standings {
		if st.PointsFor > maxPF {
			maxPF = st.PointsFor
		}
	}
	if maxPF == 0 {
		maxPF = 1 // avoid division by zero
	}

	// Fetch matchup data for streak calculation
	leagueID, err := s.sleeperSvc.GetLeagueID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get league id: %w", err)
	}

	state, err := s.sleeperClient.GetNFLState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get nfl state: %w", err)
	}

	// Build streak data by checking recent weeks
	streaks := s.calculateStreaks(ctx, leagueID, state.Week)

	var rankings []domain.PowerRanking
	for _, st := range standings {
		totalGames := st.Wins + st.Losses + st.Ties
		winPct := 0.0
		if totalGames > 0 {
			winPct = float64(st.Wins) / float64(totalGames)
		}

		normalizedPF := st.PointsFor / maxPF
		streak := streaks[st.RosterID]
		streakBonus := calculateStreakBonus(streak)

		score := (winPct * 60) + (normalizedPF * 30) + (streakBonus * 10)
		score = math.Round(score*100) / 100

		rankings = append(rankings, domain.PowerRanking{
			OwnerID:   st.OwnerID,
			OwnerName: st.OwnerName,
			Score:     score,
			Wins:      st.Wins,
			Losses:    st.Losses,
			PointsFor: st.PointsFor,
			Streak:    formatStreak(streak),
			WinPct:    math.Round(winPct*1000) / 1000,
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings, nil
}

// calculateStreaks looks at recent matchup results per roster to compute streaks.
// Returns map[rosterID]streak where positive = win streak, negative = loss streak.
func (s *AnalyticsService) calculateStreaks(ctx context.Context, leagueID string, currentWeek int) map[int]int {
	streaks := make(map[int]int)
	determined := make(map[int]bool) // rosters whose streak has been broken or finalized

	// Walk backwards from the previous week to compute streaks
	for week := currentWeek - 1; week >= 1; week-- {
		matchups, err := s.sleeperClient.GetMatchups(ctx, leagueID, week)
		if err != nil {
			break
		}

		// Group by matchup_id
		grouped := make(map[int][]sleeper.SleeperMatchup)
		for _, m := range matchups {
			grouped[m.MatchupID] = append(grouped[m.MatchupID], m)
		}

		allDetermined := true
		for _, entries := range grouped {
			if len(entries) != 2 {
				continue
			}
			for i := 0; i < 2; i++ {
				roster := entries[i]
				opponent := entries[1-i]

				if determined[roster.RosterID] {
					continue
				}

				won := roster.Points > opponent.Points
				lost := roster.Points < opponent.Points

				if week == currentWeek-1 {
					// Initialize streak
					if won {
						streaks[roster.RosterID] = 1
					} else if lost {
						streaks[roster.RosterID] = -1
					} else {
						// Tie on most recent week — no streak
						determined[roster.RosterID] = true
					}
				} else {
					// Extend streak if same direction
					prev := streaks[roster.RosterID]
					if prev > 0 && won {
						streaks[roster.RosterID]++
					} else if prev < 0 && lost {
						streaks[roster.RosterID]--
					} else {
						// Streak broken
						determined[roster.RosterID] = true
					}
				}

				if !determined[roster.RosterID] {
					allDetermined = false
				}
			}
		}

		if allDetermined {
			break
		}

		// Stop searching after 5 weeks back
		if currentWeek-week >= 5 {
			break
		}
	}

	return streaks
}

func calculateStreakBonus(streak int) float64 {
	if streak >= 3 {
		return 1.0
	}
	if streak == 2 {
		return 0.7
	}
	if streak == 1 {
		return 0.4
	}
	if streak == -1 {
		return -0.2
	}
	if streak == -2 {
		return -0.4
	}
	if streak <= -3 {
		return -0.6
	}
	return 0.0
}

func formatStreak(streak int) string {
	if streak > 0 {
		return fmt.Sprintf("W%d", streak)
	}
	if streak < 0 {
		return fmt.Sprintf("L%d", -streak)
	}
	return "-"
}

// GetAwardProjections returns projected award races based on current standings.
func (s *AnalyticsService) GetAwardProjections(ctx context.Context) ([]domain.AwardProjection, error) {
	standings, err := s.sleeperSvc.GetLiveStandings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get live standings: %w", err)
	}

	if len(standings) == 0 {
		return nil, nil
	}

	// Juggernaut: highest points for
	pfSorted := make([]domain.LiveStanding, len(standings))
	copy(pfSorted, standings)
	sort.Slice(pfSorted, func(i, j int) bool {
		return pfSorted[i].PointsFor > pfSorted[j].PointsFor
	})

	juggernaut := domain.AwardProjection{AwardName: "Juggernaut"}
	for i, st := range pfSorted {
		if i >= 5 {
			break
		}
		juggernaut.Contenders = append(juggernaut.Contenders, domain.AwardProjectionEntry{
			OwnerID:   st.OwnerID,
			OwnerName: st.OwnerName,
			StatValue: st.PointsFor,
			Rank:      i + 1,
		})
	}

	// Cupcake: lowest points for (among active owners)
	cupcake := domain.AwardProjection{AwardName: "Cupcake"}
	for i := len(pfSorted) - 1; i >= 0 && len(cupcake.Contenders) < 5; i-- {
		st := pfSorted[i]
		cupcake.Contenders = append(cupcake.Contenders, domain.AwardProjectionEntry{
			OwnerID:   st.OwnerID,
			OwnerName: st.OwnerName,
			StatValue: st.PointsFor,
			Rank:      len(cupcake.Contenders) + 1,
		})
	}

	// Bottomfeeder: most losses
	lossSorted := make([]domain.LiveStanding, len(standings))
	copy(lossSorted, standings)
	sort.Slice(lossSorted, func(i, j int) bool {
		if lossSorted[i].Losses != lossSorted[j].Losses {
			return lossSorted[i].Losses > lossSorted[j].Losses
		}
		return lossSorted[i].PointsFor < lossSorted[j].PointsFor
	})

	bottomfeeder := domain.AwardProjection{AwardName: "Bottomfeeder"}
	for i, st := range lossSorted {
		if i >= 5 {
			break
		}
		bottomfeeder.Contenders = append(bottomfeeder.Contenders, domain.AwardProjectionEntry{
			OwnerID:   st.OwnerID,
			OwnerName: st.OwnerName,
			StatValue: float64(st.Losses),
			Rank:      i + 1,
		})
	}

	// League Champion: most wins (playoff projection)
	winSorted := make([]domain.LiveStanding, len(standings))
	copy(winSorted, standings)
	sort.Slice(winSorted, func(i, j int) bool {
		if winSorted[i].Wins != winSorted[j].Wins {
			return winSorted[i].Wins > winSorted[j].Wins
		}
		return winSorted[i].PointsFor > winSorted[j].PointsFor
	})

	champion := domain.AwardProjection{AwardName: "League Champion"}
	for i, st := range winSorted {
		if i >= 5 {
			break
		}
		champion.Contenders = append(champion.Contenders, domain.AwardProjectionEntry{
			OwnerID:   st.OwnerID,
			OwnerName: st.OwnerName,
			StatValue: float64(st.Wins),
			Rank:      i + 1,
		})
	}

	return []domain.AwardProjection{champion, juggernaut, cupcake, bottomfeeder}, nil
}
