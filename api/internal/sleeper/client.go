package sleeper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const baseURL = "https://api.sleeper.app/v1"

// Client is an HTTP client for the Sleeper Fantasy Football API
// with built-in rate limiting (< 1000 calls/min).
type Client struct {
	http    *http.Client
	limiter *rateLimiter
	cache   *Cache
}

// NewClient creates a Sleeper API client with rate limiting.
func NewClient(cache *Cache) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		limiter: newRateLimiter(900, time.Minute), // 900/min to stay under 1000
		cache:   cache,
	}
}

func (c *Client) get(ctx context.Context, path string, result any) error {
	if err := c.limiter.wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}

	url := baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sleeper request %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sleeper %s returned %d: %s", path, resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// getOrCache fetches from cache first, falling back to the API.
func (c *Client) getOrCache(ctx context.Context, cacheKey, path string, ttl time.Duration, result any) error {
	if data, ok := c.cache.Get(cacheKey); ok {
		return json.Unmarshal(data, result)
	}

	if err := c.get(ctx, path, result); err != nil {
		return err
	}

	data, err := json.Marshal(result)
	if err == nil {
		c.cache.Set(cacheKey, data, ttl)
	}
	return nil
}

// NFLState represents the current NFL season state.
type NFLState struct {
	Season       string `json:"season"`
	Week         int    `json:"week"`
	SeasonType   string `json:"season_type"`
	DisplayWeek  int    `json:"display_week"`
	Leg          int    `json:"leg"`
}

// LeagueInfo represents Sleeper league metadata.
type LeagueInfo struct {
	LeagueID    string `json:"league_id"`
	Name        string `json:"name"`
	Season      string `json:"season"`
	Status      string `json:"status"`
	TotalRosters int   `json:"total_rosters"`
	Settings    map[string]any `json:"settings"`
}

// SleeperUser represents a user in a Sleeper league.
type SleeperUser struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	TeamName    string `json:"metadata,omitempty"`
}

// SleeperRoster represents a roster in a Sleeper league.
type SleeperRoster struct {
	RosterID int               `json:"roster_id"`
	OwnerID  string            `json:"owner_id"`
	Settings RosterSettings    `json:"settings"`
	Players  []string          `json:"players"`
	Starters []string          `json:"starters"`
}

// RosterSettings contains W/L/T and points.
type RosterSettings struct {
	Wins              int `json:"wins"`
	Losses            int `json:"losses"`
	Ties              int `json:"ties"`
	Fpts              int `json:"fpts"`
	FptsDecimal       int `json:"fpts_decimal"`
	FptsAgainst       int `json:"fpts_against"`
	FptsAgainstDecimal int `json:"fpts_against_decimal"`
}

// PointsFor returns the combined points-for value (e.g., 1234 + 56 = 1234.56).
func (s RosterSettings) PointsFor() float64 {
	return float64(s.Fpts) + float64(s.FptsDecimal)/100.0
}

// PointsAgainst returns the combined points-against value.
func (s RosterSettings) PointsAgainst() float64 {
	return float64(s.FptsAgainst) + float64(s.FptsAgainstDecimal)/100.0
}

// SleeperMatchup represents a single matchup entry for a given week.
type SleeperMatchup struct {
	RosterID  int      `json:"roster_id"`
	MatchupID int      `json:"matchup_id"`
	Points    float64  `json:"points"`
	Starters  []string `json:"starters"`
}

// GetNFLState returns the current NFL season/week state.
func (c *Client) GetNFLState(ctx context.Context) (*NFLState, error) {
	var state NFLState
	cacheKey := "nfl_state"
	if err := c.getOrCache(ctx, cacheKey, "/state/nfl", 15*time.Minute, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// GetLeague returns league info for the given league ID.
func (c *Client) GetLeague(ctx context.Context, leagueID string) (*LeagueInfo, error) {
	var league LeagueInfo
	cacheKey := fmt.Sprintf("league:%s", leagueID)
	if err := c.getOrCache(ctx, cacheKey, fmt.Sprintf("/league/%s", leagueID), 1*time.Hour, &league); err != nil {
		return nil, err
	}
	return &league, nil
}

// GetUsers returns users in a league.
func (c *Client) GetUsers(ctx context.Context, leagueID string) ([]SleeperUser, error) {
	var users []SleeperUser
	cacheKey := fmt.Sprintf("users:%s", leagueID)
	if err := c.getOrCache(ctx, cacheKey, fmt.Sprintf("/league/%s/users", leagueID), 1*time.Hour, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetRosters returns rosters for a league.
func (c *Client) GetRosters(ctx context.Context, leagueID string) ([]SleeperRoster, error) {
	var rosters []SleeperRoster
	cacheKey := fmt.Sprintf("rosters:%s", leagueID)
	if err := c.getOrCache(ctx, cacheKey, fmt.Sprintf("/league/%s/rosters", leagueID), 15*time.Minute, &rosters); err != nil {
		return nil, err
	}
	return rosters, nil
}

// GetMatchups returns matchups for a specific week.
func (c *Client) GetMatchups(ctx context.Context, leagueID string, week int) ([]SleeperMatchup, error) {
	var matchups []SleeperMatchup
	cacheKey := fmt.Sprintf("matchups:%s:%d", leagueID, week)
	if err := c.getOrCache(ctx, cacheKey, fmt.Sprintf("/league/%s/matchups/%d", leagueID, week), 30*time.Second, &matchups); err != nil {
		return nil, err
	}
	return matchups, nil
}

// rateLimiter implements a sliding window rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	timestamps []time.Time
	maxCalls   int
	window     time.Duration
}

func newRateLimiter(maxCalls int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		timestamps: make([]time.Time, 0, maxCalls),
		maxCalls:   maxCalls,
		window:     window,
	}
}

func (rl *rateLimiter) wait(ctx context.Context) error {
	for {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Remove expired timestamps
		valid := rl.timestamps[:0]
		for _, t := range rl.timestamps {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}
		rl.timestamps = valid

		if len(rl.timestamps) < rl.maxCalls {
			rl.timestamps = append(rl.timestamps, now)
			rl.mu.Unlock()
			return nil
		}

		// Calculate wait time until oldest entry expires
		waitUntil := rl.timestamps[0].Add(rl.window)
		rl.mu.Unlock()

		wait := time.Until(waitUntil)
		log.Debug().Dur("wait", wait).Msg("rate limit: waiting")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}
