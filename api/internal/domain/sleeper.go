package domain

import "time"

// OwnerSleeperMapping maps a local owner to their Sleeper user ID.
type OwnerSleeperMapping struct {
	ID            int       `json:"id"`
	OwnerID       int       `json:"owner_id"`
	SleeperUserID string    `json:"sleeper_user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// OwnerSleeperMappingWithName includes the owner's name for API responses.
type OwnerSleeperMappingWithName struct {
	ID            int    `json:"id"`
	OwnerID       int    `json:"owner_id"`
	OwnerName     string `json:"owner_name"`
	SleeperUserID string `json:"sleeper_user_id"`
}

// SleeperLeagueConfig maps a season to its Sleeper league ID.
type SleeperLeagueConfig struct {
	ID              int       `json:"id"`
	SeasonID        int       `json:"season_id"`
	SleeperLeagueID string    `json:"sleeper_league_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// CacheEntry represents a persistent cache row.
type CacheEntry struct {
	ID        int       `json:"id"`
	CacheKey  string    `json:"cache_key"`
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

// LiveMatchup represents a formatted live matchup for the API response.
type LiveMatchup struct {
	MatchupID int              `json:"matchup_id"`
	Teams     []LiveMatchupTeam `json:"teams"`
}

// LiveMatchupTeam represents one team in a live matchup.
type LiveMatchupTeam struct {
	OwnerID   *int    `json:"owner_id,omitempty"`
	OwnerName string  `json:"owner_name"`
	RosterID  int     `json:"roster_id"`
	Points    float64 `json:"points"`
}

// LiveStanding represents a team's current standing from Sleeper.
type LiveStanding struct {
	OwnerID   *int    `json:"owner_id,omitempty"`
	OwnerName string  `json:"owner_name"`
	RosterID  int     `json:"roster_id"`
	Wins      int     `json:"wins"`
	Losses    int     `json:"losses"`
	Ties      int     `json:"ties"`
	PointsFor float64 `json:"points_for"`
}

// LiveBeltStatus represents the current belt holder with their live matchup.
type LiveBeltStatus struct {
	CurrentHolder  string          `json:"current_holder"`
	Week           int             `json:"week"`
	Season         string          `json:"season"`
	CurrentMatchup *LiveMatchup    `json:"current_matchup,omitempty"`
}
