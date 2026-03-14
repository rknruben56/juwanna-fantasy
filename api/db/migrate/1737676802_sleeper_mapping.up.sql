-- Sleeper integration tables

-- Maps local owners to Sleeper user IDs
CREATE TABLE owner_sleeper_mapping (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    sleeper_user_id VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_owner_sleeper UNIQUE (owner_id)
);

COMMENT ON TABLE owner_sleeper_mapping IS 'Maps local owner IDs to Sleeper platform user IDs';

-- Stores Sleeper league IDs per season
CREATE TABLE sleeper_league_config (
    id SERIAL PRIMARY KEY,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    sleeper_league_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_league_season UNIQUE (season_id)
);

COMMENT ON TABLE sleeper_league_config IS 'Sleeper league IDs for each season';

-- Caches Sleeper API responses
CREATE TABLE sleeper_cache (
    id SERIAL PRIMARY KEY,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    data JSONB NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE sleeper_cache IS 'Persistent cache for Sleeper API responses';
CREATE INDEX idx_sleeper_cache_key ON sleeper_cache(cache_key);
CREATE INDEX idx_sleeper_cache_expires ON sleeper_cache(expires_at);

-- Pre-computed head-to-head records
CREATE TABLE head_to_head_records (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    opponent_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    total_points_for DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_points_against DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_h2h_pair UNIQUE (owner_id, opponent_id),
    CONSTRAINT no_self_h2h CHECK (owner_id != opponent_id)
);

CREATE INDEX idx_h2h_owner ON head_to_head_records(owner_id);
CREATE INDEX idx_h2h_opponent ON head_to_head_records(opponent_id);

-- Historical weekly matchup results
CREATE TABLE matchup_history (
    id SERIAL PRIMARY KEY,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    week_number SMALLINT NOT NULL CHECK (week_number >= 1 AND week_number <= 18),
    owner_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    opponent_id INTEGER NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    points_scored DECIMAL(10,2) NOT NULL,
    opponent_points DECIMAL(10,2) NOT NULL,
    is_playoff BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_matchup UNIQUE (season_id, week_number, owner_id)
);

CREATE INDEX idx_matchup_owner ON matchup_history(owner_id);
CREATE INDEX idx_matchup_season ON matchup_history(season_id);

-- View: Aggregated rivalry statistics
CREATE OR REPLACE VIEW rivalry_stats AS
SELECT
    h.owner_id,
    o1.name AS owner_name,
    h.opponent_id,
    o2.name AS opponent_name,
    h.wins,
    h.losses,
    h.wins + h.losses AS total_games,
    CASE WHEN h.wins + h.losses > 0
        THEN ROUND(h.wins::DECIMAL / (h.wins + h.losses), 3)
        ELSE 0
    END AS win_pct,
    h.total_points_for,
    h.total_points_against
FROM head_to_head_records h
JOIN owners o1 ON h.owner_id = o1.id
JOIN owners o2 ON h.opponent_id = o2.id
ORDER BY total_games DESC, win_pct DESC;

COMMENT ON VIEW rivalry_stats IS 'Aggregated rivalry statistics between owners';
