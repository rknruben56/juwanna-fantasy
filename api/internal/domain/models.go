package domain

import "time"

type Owner struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	JoinYear  int        `json:"join_year"`
	LeaveYear *int       `json:"leave_year,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Season struct {
	ID        int       `json:"id"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}

type AwardType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type SeasonRecord struct {
	ID                   int      `json:"id"`
	OwnerID              int      `json:"owner_id"`
	SeasonID             int      `json:"season_id"`
	RegularSeasonWins    int      `json:"regular_season_wins"`
	RegularSeasonLosses  int      `json:"regular_season_losses"`
	PlayoffWins          int      `json:"playoff_wins"`
	PlayoffLosses        int      `json:"playoff_losses"`
	PointsFor            *float64 `json:"points_for,omitempty"`
	PointsAgainst        *float64 `json:"points_against,omitempty"`
	MadePlayoffs         bool     `json:"made_playoffs"`
	FinalStanding        *int     `json:"final_standing,omitempty"`
	TeamName             *string  `json:"team_name,omitempty"`
}

type SeasonAward struct {
	ID          int      `json:"id"`
	OwnerID     int      `json:"owner_id"`
	SeasonID    int      `json:"season_id"`
	AwardTypeID int      `json:"award_type_id"`
	StatValue   *float64 `json:"stat_value,omitempty"`
	TeamName    *string  `json:"team_name,omitempty"`
}

type BeltHistory struct {
	ID         int `json:"id"`
	OwnerID    int `json:"owner_id"`
	SeasonID   int `json:"season_id"`
	WeekNumber int `json:"week_number"`
}

// View models

type OwnerCareerStats struct {
	OwnerID                  int    `json:"owner_id"`
	Name                     string `json:"name"`
	JoinYear                 int    `json:"join_year"`
	LeaveYear                *int   `json:"leave_year,omitempty"`
	YearsInLeague            int    `json:"years_in_league"`
	TotalRegularSeasonWins   int    `json:"total_regular_season_wins"`
	TotalRegularSeasonLosses int    `json:"total_regular_season_losses"`
	TotalPlayoffWins         int    `json:"total_playoff_wins"`
	TotalPlayoffLosses       int    `json:"total_playoff_losses"`
	PlayoffAppearances       int    `json:"playoff_appearances"`
	Championships            int    `json:"championships"`
	FirstPlaceLoserTrophies  int    `json:"first_place_loser_trophies"`
}

type OwnerBeltWeeks struct {
	OwnerID       int    `json:"owner_id"`
	Name          string `json:"name"`
	WeeksWithBelt int    `json:"weeks_with_belt"`
}

type SeasonStanding struct {
	Year               int      `json:"year"`
	OwnerName          string   `json:"owner_name"`
	RegularSeasonWins  int      `json:"regular_season_wins"`
	RegularSeasonLosses int     `json:"regular_season_losses"`
	PointsFor          *float64 `json:"points_for,omitempty"`
	PointsAgainst      *float64 `json:"points_against,omitempty"`
	MadePlayoffs       bool     `json:"made_playoffs"`
	TeamName           *string  `json:"team_name,omitempty"`
}

// Response types for enriched API responses

type OwnerWithStats struct {
	Owner
	CareerStats *OwnerCareerStats `json:"career_stats,omitempty"`
}

type SeasonAwardWithDetails struct {
	ID        int      `json:"id"`
	OwnerName string   `json:"owner_name"`
	AwardName string   `json:"award_name"`
	Year      int      `json:"year"`
	StatValue *float64 `json:"stat_value,omitempty"`
	TeamName  *string  `json:"team_name,omitempty"`
}

type BeltHistoryWithDetails struct {
	OwnerName  string `json:"owner_name"`
	Year       int    `json:"year"`
	WeekNumber int    `json:"week_number"`
}

// Head-to-head and rivalry types

type HeadToHeadRecord struct {
	OwnerID            int     `json:"owner_id"`
	OwnerName          string  `json:"owner_name"`
	OpponentID         int     `json:"opponent_id"`
	OpponentName       string  `json:"opponent_name"`
	Wins               int     `json:"wins"`
	Losses             int     `json:"losses"`
	TotalGames         int     `json:"total_games"`
	WinPct             float64 `json:"win_pct"`
	TotalPointsFor     float64 `json:"total_points_for"`
	TotalPointsAgainst float64 `json:"total_points_against"`
}

type MatchupHistoryEntry struct {
	Season         int     `json:"season"`
	Week           int     `json:"week"`
	OwnerName      string  `json:"owner_name"`
	OpponentName   string  `json:"opponent_name"`
	PointsScored   float64 `json:"points_scored"`
	OpponentPoints float64 `json:"opponent_points"`
	IsPlayoff      bool    `json:"is_playoff"`
	Result         string  `json:"result"` // "W", "L", "T"
}

type RivalrySummary struct {
	HeadToHeadRecord
	RecentMatchups []MatchupHistoryEntry `json:"recent_matchups"`
}

// Power rankings types

type PowerRanking struct {
	Rank            int     `json:"rank"`
	OwnerID         *int    `json:"owner_id,omitempty"`
	OwnerName       string  `json:"owner_name"`
	Score           float64 `json:"score"`
	Wins            int     `json:"wins"`
	Losses          int     `json:"losses"`
	PointsFor       float64 `json:"points_for"`
	Streak          string  `json:"streak"`
	WinPct          float64 `json:"win_pct"`
}

// Award projection types

type AwardProjection struct {
	AwardName   string                   `json:"award_name"`
	Contenders  []AwardProjectionEntry   `json:"contenders"`
}

type AwardProjectionEntry struct {
	OwnerID   *int    `json:"owner_id,omitempty"`
	OwnerName string  `json:"owner_name"`
	StatValue float64 `json:"stat_value"`
	Rank      int     `json:"rank"`
}
