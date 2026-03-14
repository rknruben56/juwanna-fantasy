# Juwanna Fantasy Go API Implementation Plan

## Overview

Build a Go REST API that combines historical league data (PostgreSQL) with live Sleeper Fantasy Football API data. The API runs alongside the existing PostgreSQL container in docker-compose.

## Project Structure

```
api/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/config.go            # Environment config
│   ├── domain/                     # Entity types & interfaces
│   ├── repository/postgres/        # DB implementations
│   ├── service/                    # Business logic
│   ├── sleeper/                    # Sleeper API client + cache
│   └── handler/                    # HTTP handlers
├── db/migrate/                     # SQL migrations (existing + new)
├── Dockerfile                      # Go container
├── docker-compose.yml              # Updated with api service
├── go.mod
└── .env.example
```

## Database Schema Additions

New migration `1737676802_sleeper_mapping.up.sql`:

| Table | Purpose |
|-------|---------|
| `owner_sleeper_mapping` | Maps local owner_id to sleeper_user_id (manual maintenance) |
| `sleeper_league_config` | Stores Sleeper league IDs per season |
| `sleeper_cache` | Persists cached Sleeper API responses (JSONB) |
| `head_to_head_records` | Pre-computed H2H records between owners |
| `matchup_history` | Historical weekly matchup results |
| `rivalry_stats` (view) | Aggregated rivalry statistics |

## API Endpoints

### Historical Data
- `GET /api/v1/owners` - List owners (filter: `?active=true`)
- `GET /api/v1/owners/:id/stats` - Career stats
- `GET /api/v1/owners/:id/vs/:opponent_id` - Head-to-head record
- `GET /api/v1/seasons/:year/standings` - Season leaderboard
- `GET /api/v1/belt/current` - Current belt holder
- `GET /api/v1/belt/leaderboard` - Weeks held rankings

### Live (Sleeper-powered)
- `GET /api/v1/live/matchups` - Current week matchups
- `GET /api/v1/live/standings` - Current standings
- `GET /api/v1/live/belt` - Belt status with current matchup

### Analytics
- `GET /api/v1/rankings/power` - Power rankings
- `GET /api/v1/projections/awards` - Award race projections
- `GET /api/v1/digest/weekly` - Combined weekly summary

### Admin
- `GET/POST/PUT/DELETE /api/v1/admin/mappings` - Owner-to-Sleeper mapping CRUD

## Dependencies

```go
github.com/go-chi/chi/v5       // Router
github.com/jackc/pgx/v5        // PostgreSQL driver
github.com/kelseyhightower/envconfig  // Config
github.com/rs/zerolog          // Logging
```

## Sleeper Integration

**Caching Strategy:**
| Data | TTL |
|------|-----|
| Matchups (live) | 30 seconds |
| Rosters | 15 minutes |
| League info | 1 hour |
| Players DB | 24 hours |

**Owner Resolution Flow:**
```
Sleeper roster_id -> sleeper_user_id -> owner_sleeper_mapping -> local owner_id
```

## Implementation Phases

### Phase 1: Foundation
- [ ] Initialize Go module with dependencies
- [ ] Create Dockerfile
- [ ] Update docker-compose.yml with api service
- [ ] Implement config loading
- [ ] Set up database connection (pgx pool)
- [ ] Create router with middleware (logging, CORS)
- [ ] Health check endpoint

### Phase 2: Historical Data API
- [ ] Domain types (Owner, Season, Award, Record, Belt)
- [ ] Repository interfaces + PostgreSQL implementations
- [ ] Service layer
- [ ] Handlers for owners, seasons, awards, belt endpoints
- [ ] Add new database migration for Sleeper tables

### Phase 3: Sleeper Integration
- [ ] HTTP client with rate limiting (< 1000 calls/min)
- [ ] In-memory cache with TTL
- [ ] Database cache persistence
- [ ] Owner mapping admin endpoints
- [ ] Live endpoints (matchups, standings, rosters)

### Phase 4: Rivalries & Analytics
- [ ] Head-to-head calculation logic
- [ ] Rivalry endpoints
- [ ] Power rankings algorithm
- [ ] Award race projections

### Phase 5: Weekly Digest & Polish
- [ ] Digest service (matchups + H2H context + belt + awards)
- [ ] API documentation
- [ ] Production config

## Docker Setup

**Updated docker-compose.yml:**
```yaml
services:
  postgres:
    # ... existing config ...
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_HOST=postgres
      - SLEEPER_LEAGUE_ID=${SLEEPER_LEAGUE_ID}
    depends_on:
      postgres:
        condition: service_healthy
```

## Files to Create/Modify

| File | Action |
|------|--------|
| `api/Dockerfile` | Create |
| `api/docker-compose.yml` | Modify (add api service) |
| `api/go.mod` | Create |
| `api/.env.example` | Create |
| `api/cmd/server/main.go` | Create |
| `api/internal/**/*.go` | Create (all packages) |
| `api/db/migrate/1737676802_sleeper_mapping.up.sql` | Create |
| `api/db/migrate/1737676802_sleeper_mapping.down.sql` | Create |
| `CLAUDE.md` | Update with new commands |

## Verification

1. `docker compose up --build` - Both containers start
2. `curl localhost:8080/health` - API responds
3. `curl localhost:8080/api/v1/owners` - Returns owners from DB
4. `curl localhost:8080/api/v1/live/matchups` - Returns Sleeper data (after mapping configured)
