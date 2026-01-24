-- Rollback initial schema migration

DROP VIEW IF EXISTS season_leaderboard;
DROP VIEW IF EXISTS owner_belt_weeks;
DROP VIEW IF EXISTS owner_career_stats;

DROP TABLE IF EXISTS championship_belt_history CASCADE;
DROP TABLE IF EXISTS season_awards CASCADE;
DROP TABLE IF EXISTS season_records CASCADE;
DROP TABLE IF EXISTS award_types CASCADE;
DROP TABLE IF EXISTS seasons CASCADE;
DROP TABLE IF EXISTS owners CASCADE;
