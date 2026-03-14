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
