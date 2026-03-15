# Juwanna Fantasy API Documentation

Base URL: `http://localhost:8080`

All responses are JSON. Errors return `{"error": "message"}`.

---

## Health

### `GET /health`

Returns API and database status.

**Response:**
```json
{"status": "ok", "db": "up"}
```

---

## Historical Data

### `GET /api/v1/owners`

List all owners. Filter with `?active=true` for active owners only.

**Query Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `active` | string | Set to `"true"` to filter active owners |

**Response:** `Owner[]`
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "join_year": 2015,
    "leave_year": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### `GET /api/v1/owners/{id}/stats`

Get career statistics for an owner.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `id` | int | Owner ID |

**Response:** `OwnerCareerStats`
```json
{
  "owner_id": 1,
  "name": "John Doe",
  "join_year": 2015,
  "years_in_league": 11,
  "total_regular_season_wins": 85,
  "total_regular_season_losses": 63,
  "total_playoff_wins": 8,
  "total_playoff_losses": 5,
  "playoff_appearances": 7,
  "championships": 2,
  "first_place_loser_trophies": 1
}
```

### `GET /api/v1/seasons/{year}/standings`

Get standings for a season.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `year` | int | Season year (e.g., 2025) |

**Response:** `SeasonStanding[]`
```json
[
  {
    "year": 2025,
    "owner_name": "John Doe",
    "regular_season_wins": 10,
    "regular_season_losses": 4,
    "points_for": 1850.5,
    "points_against": 1620.3,
    "made_playoffs": true,
    "team_name": "The Destroyers"
  }
]
```

### `GET /api/v1/seasons/{year}/awards`

Get awards for a season.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `year` | int | Season year |

**Response:** `SeasonAwardWithDetails[]`
```json
[
  {
    "id": 1,
    "owner_name": "John Doe",
    "award_name": "League Champion",
    "year": 2025,
    "stat_value": null,
    "team_name": "The Destroyers"
  }
]
```

### `GET /api/v1/belt/current`

Get the current championship belt holder.

**Response:** `BeltHistoryWithDetails`
```json
{
  "owner_name": "John Doe",
  "year": 2025,
  "week_number": 12
}
```

### `GET /api/v1/belt/leaderboard`

Get belt weeks held rankings.

**Response:** `OwnerBeltWeeks[]`
```json
[
  {
    "owner_id": 1,
    "name": "John Doe",
    "weeks_with_belt": 28
  }
]
```

---

## Live (Sleeper-powered)

These endpoints pull live data from the Sleeper Fantasy Football API. Requires `SLEEPER_LEAGUE_ID` to be configured.

### `GET /api/v1/live/matchups`

Get current week matchups.

**Response:** `LiveMatchup[]`
```json
[
  {
    "matchup_id": 1,
    "teams": [
      {
        "owner_id": 1,
        "owner_name": "John Doe",
        "roster_id": 1,
        "points": 105.5
      },
      {
        "owner_id": 2,
        "owner_name": "Jane Smith",
        "roster_id": 5,
        "points": 98.2
      }
    ]
  }
]
```

### `GET /api/v1/live/standings`

Get current season standings.

**Response:** `LiveStanding[]`
```json
[
  {
    "owner_id": 1,
    "owner_name": "John Doe",
    "roster_id": 1,
    "wins": 8,
    "losses": 4,
    "ties": 0,
    "points_for": 1520.8
  }
]
```

### `GET /api/v1/live/belt`

Get belt status with current matchup details.

**Response:** `LiveBeltStatus`
```json
{
  "current_holder": "John Doe",
  "week": 13,
  "season": "2025",
  "current_matchup": {
    "matchup_id": 3,
    "teams": [...]
  }
}
```

---

## Analytics

### `GET /api/v1/owners/{id}/vs/{opponent_id}`

Get head-to-head record between two owners with recent matchup history.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `id` | int | Owner ID |
| `opponent_id` | int | Opponent owner ID |

**Response:** `RivalrySummary`
```json
{
  "owner_id": 1,
  "owner_name": "John Doe",
  "opponent_id": 2,
  "opponent_name": "Jane Smith",
  "wins": 12,
  "losses": 8,
  "total_games": 20,
  "win_pct": 0.6,
  "total_points_for": 2100.5,
  "total_points_against": 1980.3,
  "recent_matchups": [
    {
      "season": 2025,
      "week": 8,
      "owner_name": "John Doe",
      "opponent_name": "Jane Smith",
      "points_scored": 125.4,
      "opponent_points": 110.2,
      "is_playoff": false,
      "result": "W"
    }
  ]
}
```

### `GET /api/v1/owners/{id}/rivalries`

Get all head-to-head records for an owner.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `id` | int | Owner ID |

**Response:** `HeadToHeadRecord[]`

### `GET /api/v1/rankings/power`

Get power rankings computed from live Sleeper data.

Score = (WinPct * 60) + (NormalizedPointsFor * 30) + (StreakBonus * 10)

**Response:** `PowerRanking[]`
```json
[
  {
    "rank": 1,
    "owner_id": 1,
    "owner_name": "John Doe",
    "score": 78.5,
    "wins": 10,
    "losses": 3,
    "points_for": 1650.2,
    "streak": "W3",
    "win_pct": 0.769
  }
]
```

### `GET /api/v1/projections/awards`

Get projected award races based on current standings.

**Response:** `AwardProjection[]`
```json
[
  {
    "award_name": "League Champion",
    "contenders": [
      {
        "owner_id": 1,
        "owner_name": "John Doe",
        "stat_value": 10,
        "rank": 1
      }
    ]
  }
]
```

Awards projected: League Champion, Juggernaut, Cupcake, Bottomfeeder.

---

## Digest

### `GET /api/v1/digest/weekly`

Combined weekly summary with matchups (enriched with H2H context), standings, belt status, award projections, and power rankings.

**Response:** `WeeklyDigest`
```json
{
  "season": "2025",
  "week": 13,
  "matchups": [
    {
      "matchup_id": 1,
      "teams": [...],
      "head_to_head": {
        "owner_wins": 12,
        "opponent_wins": 8,
        "total_games": 20
      }
    }
  ],
  "standings": [...],
  "belt_status": {...},
  "award_projections": [...],
  "power_rankings": [...]
}
```

---

## Admin

### `GET /api/v1/admin/mappings`

List all owner-to-Sleeper user ID mappings.

**Response:** `OwnerSleeperMappingWithName[]`
```json
[
  {
    "id": 1,
    "owner_id": 1,
    "owner_name": "John Doe",
    "sleeper_user_id": "123456789"
  }
]
```

### `POST /api/v1/admin/mappings`

Create a new owner-to-Sleeper mapping.

**Request Body:**
```json
{
  "owner_id": 1,
  "sleeper_user_id": "123456789"
}
```

**Response:** `OwnerSleeperMapping` (201 Created)

### `PUT /api/v1/admin/mappings/{id}`

Update an existing mapping.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `id` | int | Mapping ID |

**Request Body:**
```json
{
  "sleeper_user_id": "987654321"
}
```

**Response:** `OwnerSleeperMapping`

### `DELETE /api/v1/admin/mappings/{id}`

Delete a mapping.

**Path Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `id` | int | Mapping ID |

**Response:** 204 No Content
