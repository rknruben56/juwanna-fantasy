package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
	"github.com/rknruben56/juwanna-fantasy/api/internal/sleeper"
)

// SleeperService handles live data from the Sleeper API.
type SleeperService struct {
	client       *sleeper.Client
	mappings     domain.MappingRepository
	leagueConfig domain.LeagueConfigRepository
	belt         domain.BeltRepository
	defaultLeagueID string
}

// NewSleeperService creates a new SleeperService.
func NewSleeperService(
	client *sleeper.Client,
	mappings domain.MappingRepository,
	leagueConfig domain.LeagueConfigRepository,
	belt domain.BeltRepository,
	defaultLeagueID string,
) *SleeperService {
	return &SleeperService{
		client:          client,
		mappings:        mappings,
		leagueConfig:    leagueConfig,
		belt:            belt,
		defaultLeagueID: defaultLeagueID,
	}
}

// GetLeagueID returns the current league ID from the database config,
// falling back to the environment variable default.
func (s *SleeperService) GetLeagueID(ctx context.Context) (string, error) {
	leagueID, err := s.leagueConfig.GetCurrentLeagueID(ctx)
	if err == nil && leagueID != "" {
		return leagueID, nil
	}
	if s.defaultLeagueID != "" {
		return s.defaultLeagueID, nil
	}
	return "", fmt.Errorf("no sleeper league ID configured")
}

// buildOwnerLookup creates a map from sleeper_user_id to owner info.
func (s *SleeperService) buildOwnerLookup(ctx context.Context) (map[string]domain.OwnerSleeperMappingWithName, error) {
	mappings, err := s.mappings.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list mappings: %w", err)
	}
	lookup := make(map[string]domain.OwnerSleeperMappingWithName, len(mappings))
	for _, m := range mappings {
		lookup[m.SleeperUserID] = m
	}
	return lookup, nil
}

// resolveOwnerName maps a Sleeper roster to a local owner name.
func resolveOwnerName(roster sleeper.SleeperRoster, users map[string]sleeper.SleeperUser, ownerLookup map[string]domain.OwnerSleeperMappingWithName) (string, *int) {
	if mapping, ok := ownerLookup[roster.OwnerID]; ok {
		ownerID := mapping.OwnerID
		return mapping.OwnerName, &ownerID
	}
	if user, ok := users[roster.OwnerID]; ok {
		return user.DisplayName, nil
	}
	return fmt.Sprintf("Roster %d", roster.RosterID), nil
}

// GetLiveMatchups returns formatted matchups for the current week.
func (s *SleeperService) GetLiveMatchups(ctx context.Context) ([]domain.LiveMatchup, error) {
	leagueID, err := s.GetLeagueID(ctx)
	if err != nil {
		return nil, err
	}

	state, err := s.client.GetNFLState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get nfl state: %w", err)
	}

	rosters, err := s.client.GetRosters(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get rosters: %w", err)
	}

	users, err := s.client.GetUsers(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	matchups, err := s.client.GetMatchups(ctx, leagueID, state.Week)
	if err != nil {
		return nil, fmt.Errorf("get matchups: %w", err)
	}

	ownerLookup, err := s.buildOwnerLookup(ctx)
	if err != nil {
		return nil, err
	}

	// Build roster lookup: roster_id -> roster
	rosterLookup := make(map[int]sleeper.SleeperRoster, len(rosters))
	for _, r := range rosters {
		rosterLookup[r.RosterID] = r
	}

	// Build user lookup: user_id -> user
	userLookup := make(map[string]sleeper.SleeperUser, len(users))
	for _, u := range users {
		userLookup[u.UserID] = u
	}

	// Group matchups by matchup_id
	grouped := make(map[int][]sleeper.SleeperMatchup)
	for _, m := range matchups {
		grouped[m.MatchupID] = append(grouped[m.MatchupID], m)
	}

	var result []domain.LiveMatchup
	for matchupID, entries := range grouped {
		lm := domain.LiveMatchup{MatchupID: matchupID}
		for _, entry := range entries {
			roster := rosterLookup[entry.RosterID]
			name, ownerID := resolveOwnerName(roster, userLookup, ownerLookup)
			lm.Teams = append(lm.Teams, domain.LiveMatchupTeam{
				OwnerID:   ownerID,
				OwnerName: name,
				RosterID:  entry.RosterID,
				Points:    entry.Points,
			})
		}
		result = append(result, lm)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].MatchupID < result[j].MatchupID
	})

	return result, nil
}

// GetLiveStandings returns current season standings from Sleeper.
func (s *SleeperService) GetLiveStandings(ctx context.Context) ([]domain.LiveStanding, error) {
	leagueID, err := s.GetLeagueID(ctx)
	if err != nil {
		return nil, err
	}

	rosters, err := s.client.GetRosters(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get rosters: %w", err)
	}

	users, err := s.client.GetUsers(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	ownerLookup, err := s.buildOwnerLookup(ctx)
	if err != nil {
		return nil, err
	}

	userLookup := make(map[string]sleeper.SleeperUser, len(users))
	for _, u := range users {
		userLookup[u.UserID] = u
	}

	var standings []domain.LiveStanding
	for _, roster := range rosters {
		name, ownerID := resolveOwnerName(roster, userLookup, ownerLookup)
		standings = append(standings, domain.LiveStanding{
			OwnerID:   ownerID,
			OwnerName: name,
			RosterID:  roster.RosterID,
			Wins:      roster.Settings.Wins,
			Losses:    roster.Settings.Losses,
			Ties:      roster.Settings.Ties,
			PointsFor: roster.Settings.PointsFor(),
		})
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Wins != standings[j].Wins {
			return standings[i].Wins > standings[j].Wins
		}
		return standings[i].PointsFor > standings[j].PointsFor
	})

	return standings, nil
}

// GetLiveBeltStatus returns the current belt holder and their active matchup.
func (s *SleeperService) GetLiveBeltStatus(ctx context.Context) (*domain.LiveBeltStatus, error) {
	beltHolder, err := s.belt.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("get belt holder: %w", err)
	}

	state, err := s.client.GetNFLState(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("could not get NFL state for belt status")
		return &domain.LiveBeltStatus{
			CurrentHolder: beltHolder.OwnerName,
			Week:          beltHolder.WeekNumber,
			Season:        fmt.Sprintf("%d", beltHolder.Year),
		}, nil
	}

	status := &domain.LiveBeltStatus{
		CurrentHolder: beltHolder.OwnerName,
		Week:          state.Week,
		Season:        state.Season,
	}

	// Try to find the belt holder's current matchup
	matchups, err := s.GetLiveMatchups(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("could not get matchups for belt status")
		return status, nil
	}

	for _, m := range matchups {
		for _, team := range m.Teams {
			if team.OwnerName == beltHolder.OwnerName {
				status.CurrentMatchup = &m
				break
			}
		}
		if status.CurrentMatchup != nil {
			break
		}
	}

	return status, nil
}

// Mapping CRUD methods

func (s *SleeperService) ListMappings(ctx context.Context) ([]domain.OwnerSleeperMappingWithName, error) {
	return s.mappings.List(ctx)
}

func (s *SleeperService) CreateMapping(ctx context.Context, ownerID int, sleeperUserID string) (*domain.OwnerSleeperMapping, error) {
	return s.mappings.Create(ctx, ownerID, sleeperUserID)
}

func (s *SleeperService) UpdateMapping(ctx context.Context, id int, sleeperUserID string) (*domain.OwnerSleeperMapping, error) {
	return s.mappings.Update(ctx, id, sleeperUserID)
}

func (s *SleeperService) DeleteMapping(ctx context.Context, id int) error {
	return s.mappings.Delete(ctx, id)
}
