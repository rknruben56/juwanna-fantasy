-- Initial schema migration for Juwanna Fantasy League

-- ============================================================
-- CORE TABLES
-- ============================================================

-- Owners table: Stores league member information
CREATE TABLE owners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    join_year SMALLINT NOT NULL,
    leave_year SMALLINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE owners IS 'League members/team owners';
COMMENT ON COLUMN owners.leave_year IS 'NULL indicates owner is still active';

-- Seasons table: Stores each season year
CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    year SMALLINT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE seasons IS 'Fantasy football seasons by year';

-- Award types: Lookup table for different awards
CREATE TABLE award_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

COMMENT ON TABLE award_types IS 'Types of awards given each season';

-- ============================================================
-- SEASON RECORDS
-- ============================================================

-- Season records: Owner's record for a specific season
CREATE TABLE season_records (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    regular_season_wins SMALLINT NOT NULL DEFAULT 0,
    regular_season_losses SMALLINT NOT NULL DEFAULT 0,
    playoff_wins SMALLINT NOT NULL DEFAULT 0,
    playoff_losses SMALLINT NOT NULL DEFAULT 0,
    points_for DECIMAL(10,2),
    points_against DECIMAL(10,2),
    made_playoffs BOOLEAN NOT NULL DEFAULT FALSE,
    final_standing SMALLINT,
    team_name VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_owner_season UNIQUE (owner_id, season_id)
);

COMMENT ON TABLE season_records IS 'Owner performance records per season';
COMMENT ON COLUMN season_records.final_standing IS '1 = champion, higher numbers = worse placement';

-- ============================================================
-- AWARDS
-- ============================================================

-- Season awards: Links owners to awards they won in specific seasons
CREATE TABLE season_awards (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    award_type_id INTEGER NOT NULL REFERENCES award_types(id) ON DELETE CASCADE,
    stat_value DECIMAL(10,2),
    team_name VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_award_per_season UNIQUE (season_id, award_type_id)
);

COMMENT ON TABLE season_awards IS 'Awards won by owners each season';
COMMENT ON COLUMN season_awards.stat_value IS 'Associated stat (e.g., PF for Juggernaut, PA for Unlucky/Cupcake)';

-- ============================================================
-- CHAMPIONSHIP BELT
-- ============================================================

-- Championship belt history: Tracks who held the belt each week
CREATE TABLE championship_belt_history (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    week_number SMALLINT NOT NULL CHECK (week_number >= 1 AND week_number <= 18),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_belt_week UNIQUE (season_id, week_number)
);

COMMENT ON TABLE championship_belt_history IS 'Weekly championship belt holder tracking';

-- ============================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================

CREATE INDEX idx_season_records_owner ON season_records(owner_id);
CREATE INDEX idx_season_records_season ON season_records(season_id);
CREATE INDEX idx_season_records_playoffs ON season_records(made_playoffs) WHERE made_playoffs = TRUE;

CREATE INDEX idx_season_awards_owner ON season_awards(owner_id);
CREATE INDEX idx_season_awards_season ON season_awards(season_id);
CREATE INDEX idx_season_awards_type ON season_awards(award_type_id);

CREATE INDEX idx_belt_history_owner ON championship_belt_history(owner_id);
CREATE INDEX idx_belt_history_season ON championship_belt_history(season_id);

-- ============================================================
-- VIEWS
-- ============================================================

-- View: Owner career statistics
CREATE OR REPLACE VIEW owner_career_stats AS
SELECT
    o.id AS owner_id,
    o.name,
    o.join_year,
    o.leave_year,
    CASE
        WHEN o.leave_year IS NULL THEN EXTRACT(YEAR FROM CURRENT_DATE)::INT - o.join_year + 1
        ELSE o.leave_year - o.join_year + 1
    END AS years_in_league,
    COALESCE(SUM(sr.regular_season_wins), 0) AS total_regular_season_wins,
    COALESCE(SUM(sr.regular_season_losses), 0) AS total_regular_season_losses,
    COALESCE(SUM(sr.playoff_wins), 0) AS total_playoff_wins,
    COALESCE(SUM(sr.playoff_losses), 0) AS total_playoff_losses,
    COUNT(sr.id) FILTER (WHERE sr.made_playoffs = TRUE) AS playoff_appearances,
    COUNT(sa.id) FILTER (WHERE at.name = 'League Champion') AS championships,
    COUNT(sa.id) FILTER (WHERE at.name = 'First Place Loser') AS first_place_loser_trophies
FROM owners o
LEFT JOIN season_records sr ON o.id = sr.owner_id
LEFT JOIN season_awards sa ON o.id = sa.owner_id
LEFT JOIN award_types at ON sa.award_type_id = at.id
GROUP BY o.id, o.name, o.join_year, o.leave_year;

COMMENT ON VIEW owner_career_stats IS 'Aggregated career statistics for each owner';

-- View: Championship belt weeks held per owner
CREATE OR REPLACE VIEW owner_belt_weeks AS
SELECT
    o.id AS owner_id,
    o.name,
    COUNT(cbh.id) AS weeks_with_belt
FROM owners o
LEFT JOIN championship_belt_history cbh ON o.id = cbh.owner_id
GROUP BY o.id, o.name;

COMMENT ON VIEW owner_belt_weeks IS 'Total weeks each owner held the championship belt';

-- View: Season leaderboard
CREATE OR REPLACE VIEW season_leaderboard AS
SELECT
    s.year,
    o.name AS owner_name,
    sr.regular_season_wins,
    sr.regular_season_losses,
    sr.points_for,
    sr.points_against,
    sr.made_playoffs,
    sr.team_name
FROM season_records sr
JOIN seasons s ON sr.season_id = s.id
JOIN owners o ON sr.owner_id = o.id
ORDER BY s.year DESC, sr.regular_season_wins DESC;

COMMENT ON VIEW season_leaderboard IS 'Season standings by year';
