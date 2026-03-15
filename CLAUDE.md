# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Juwanna Fantasy is a website for a Fantasy Football League that tracks historical league data including owners, season records, awards, and championship belt history.

## Development Commands

### Database

```bash
# Start all services (PostgreSQL, API, Web)
docker compose up -d

# Run migrations (requires pgmgr - https://github.com/rnubel/pgmgr)
cd api && pgmgr migrate

# Rollback migrations
cd api && pgmgr rollback
```

Database connection: `localhost:5432`, database `juwanna_fantasy`, user/password `postgres/postgres`

## Architecture

### Database Schema (PostgreSQL 16)

Core tables:
- `owners` - League members with join/leave years
- `seasons` - Season years
- `award_types` - Lookup table for awards (League Champion, First Place Loser, Unlucky, Juggernaut, Cupcake, Bottomfeeder)
- `season_records` - Per-owner per-season stats (wins, losses, points, playoff status)
- `season_awards` - Awards won per season
- `championship_belt_history` - Weekly belt holder tracking

Views for aggregated data:
- `owner_career_stats` - Career totals per owner
- `owner_belt_weeks` - Total weeks each owner held the belt
- `season_leaderboard` - Season standings

### Migration Files

Located in `api/db/migrate/` using timestamp-prefixed naming (`{timestamp}_{name}.up.sql` / `.down.sql`). Managed via pgmgr with config in `api/.pgmgr.json`.
