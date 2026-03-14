# Juwanna Fantasy Frontend UI Plan

## Context

The Go API has 6 historical data endpoints fully implemented (owners, seasons, belt). There is no frontend yet. This plan designs a premium, dark-themed fantasy football league site that makes 11 seasons of league history feel alive — built for trash talk, bragging rights, and settling arguments.

## Design Language

**Dark theme only** — fantasy football is a nighttime sport.

| Token | Value | Usage |
|-------|-------|-------|
| Navy | `#0B1120` | Page background |
| Gold | `#D4AF37` | Champions, belt, trophies |
| Red | `#E63946` | Losses, bottomfeeder |
| Teal | `#2EC4B6` | Wins, playoffs, positive |
| Surface | `#151D2E` | Card backgrounds |
| Surface hover | `#1E293B` | Hover / secondary |
| Text | `#F1F5F9` | Primary text |
| Text muted | `#94A3B8` | Captions, secondary |

**Typography:** `Oswald` (headlines, scoreboard feel) + `Inter` (body, tables). No CSS framework — custom CSS with variables for a distinctive look.

**Award Icons:**
- League Champion → Crown (gold)
- First Place Loser → Silver medal
- Unlucky → Storm cloud (purple)
- Juggernaut → Rocket (red/orange)
- Cupcake → Cupcake (pink)
- Bottomfeeder → Dumpster fire (brown)

---

## Pages

### 1. Home Dashboard (`/`)

The league's living room. First thing members see.

| Section | Description | Data Source |
|---------|-------------|------------|
| **Belt Hero** | Full-width banner, current holder with pulsing gold glow, "Held since Week X, YYYY" | `GET /belt/current` |
| **Quick Stats Bar** | 4 stat cards: seasons, all-time owners, active owners, total belt weeks | `GET /owners`, `GET /belt/leaderboard` |
| **Belt Top 5** | Horizontal gold bar chart of top belt holders | `GET /belt/leaderboard` |
| **Latest Season** | Mini standings table + award row for most recent season | `GET /seasons/2025/standings`, `/awards` |
| **Active Owners Grid** | Cards with name, tenure, W-L, championships — link to profiles | `GET /owners?active=true`, `/owners/:id/stats` |

### 2. Owner Pages (`/owners`, `/owners/:id`)

**List:** Active/All toggle, grid of owner cards sorted by career wins. Each card: name, join year, W-L, championship count, belt weeks.

**Profile detail:**
| Section | Description |
|---------|-------------|
| **Header** | Name, tenure, row of trophy icons per championship |
| **Career Stats Grid** | 6-8 stat boxes: W-L (with %), playoff W-L, playoff apps, championships, "Bridesmaid Awards" (first place loser), belt weeks |
| **Season-by-Season Table** | Year, team name, W-L, PF, PA, playoffs, awards — built by aggregating season endpoints |
| **Award Showcase** | Visual grid of awards won, stacked for repeats ("League Champion x2") |

### 3. Season Archives (`/seasons`, `/seasons/:year`)

**List:** All 11 seasons (2015-2025) as cards with year, champion name, and mini standings preview.

**Detail:**
| Section | Description |
|---------|-------------|
| **Header** | Year + champion gold badge |
| **Standings Table** | Sortable: rank, owner, team, W-L, PF, PA, differential, playoffs. Playoff rows get green border, champion gets gold, last place gets red |
| **Award Grid** | 2x3 grid of award cards with icons. Champion card is oversized/gold. Bottomfeeder gets humorous treatment |

### 4. Championship Belt (`/belt`)

The crown jewel page. WWE title history meets stock ticker.

| Section | Description |
|---------|-------------|
| **Belt Hero (expanded)** | Enlarged belt graphic, current holder, gold glow |
| **Full Leaderboard** | Table + horizontal gold bar chart: owner, total weeks, % of all weeks. Current holder gets "CURRENT" badge |
| **Belt Timeline** | Season-by-season heatmap (18 week slots per season, color-coded by holder). *Placeholder until `GET /belt/history` endpoint is added* |
| **Fun Facts** | Computed: "Most dominant", "Longest drought", etc. |

### 5. Hall of Fame / Awards (`/awards`)

| Section | Description |
|---------|-------------|
| **Per-Award Sections** | One section per award type, vertical timeline of all winners by year |
| **Award Leaderboard** | Who has won each award the most? Most decorated overall? |
| **Superlatives** | Highest single-season PF (Juggernaut), most championships, most Bridesmaid awards, most Bottomfeeder awards |

Data assembled by calling `GET /seasons/:year/awards` for each year.

---

## Tech Stack

- **React 19 + TypeScript + Vite 5** (Vite 5 for Node 22.3 compat)
- **React Router v7** — client-side routing
- **TanStack Query v5** — server state, caching, background refetch
- **Custom CSS** with CSS variables — no framework, premium dark aesthetic
- **Google Fonts** — Oswald + Inter

---

## File Structure

```
web/src/
├── main.tsx
├── App.tsx                         # Router + QueryClientProvider
├── index.css                       # CSS variables, fonts, global reset
├── api/
│   ├── client.ts                   # Fetch wrapper
│   ├── owners.ts
│   ├── seasons.ts
│   └── belt.ts
├── hooks/
│   ├── useOwners.ts                # useOwners(), useOwnerStats(id)
│   ├── useSeasons.ts               # useSeasonStandings(year), useSeasonAwards(year)
│   └── useBelt.ts                  # useBeltCurrent(), useBeltLeaderboard()
├── types/
│   └── index.ts                    # TS interfaces mirroring Go domain models
├── components/
│   ├── layout/
│   │   ├── Navbar.tsx              # Logo, nav links, mini belt badge
│   │   ├── Footer.tsx              # "Est. 2015"
│   │   └── Layout.tsx
│   ├── shared/
│   │   ├── StatCard.tsx
│   │   ├── DataTable.tsx           # Generic sortable table
│   │   ├── Badge.tsx
│   │   ├── AwardIcon.tsx           # Award type → icon mapping
│   │   ├── LoadingSpinner.tsx
│   │   └── ErrorState.tsx
│   ├── home/
│   │   ├── BeltHero.tsx
│   │   ├── QuickStatsBar.tsx
│   │   ├── BeltLeaderboardPreview.tsx
│   │   ├── SeasonSnapshot.tsx
│   │   └── OwnerGrid.tsx
│   ├── owners/
│   │   ├── OwnerCard.tsx
│   │   ├── OwnerHeader.tsx
│   │   ├── CareerStatsGrid.tsx
│   │   ├── SeasonHistoryTable.tsx
│   │   └── AwardShowcase.tsx
│   ├── seasons/
│   │   ├── SeasonCard.tsx
│   │   ├── StandingsTable.tsx
│   │   ├── AwardGrid.tsx
│   │   └── AwardCard.tsx
│   ├── belt/
│   │   ├── BeltHeroExpanded.tsx
│   │   ├── BeltLeaderboardFull.tsx
│   │   ├── BeltTimeline.tsx        # Placeholder for future API
│   │   └── BeltFunFacts.tsx
│   └── awards/
│       ├── AwardSection.tsx
│       ├── AwardLeaderboard.tsx
│       └── SuperlativesPanel.tsx
└── pages/
    ├── HomePage.tsx
    ├── OwnersPage.tsx
    ├── OwnerProfilePage.tsx
    ├── SeasonsPage.tsx
    ├── SeasonDetailPage.tsx
    ├── BeltPage.tsx
    ├── AwardsPage.tsx
    └── NotFoundPage.tsx
```

---

## Implementation Phases

1. **Scaffold** — Vite project, deps, CSS variables, layout, router, API client, types, hooks
2. **Shared Components** — StatCard, DataTable, Badge, AwardIcon, loading/error states
3. **Home Page** — BeltHero, QuickStatsBar, BeltLeaderboardPreview, SeasonSnapshot, OwnerGrid
4. **Season Pages** — SeasonList, SeasonDetail with standings + awards
5. **Owner Pages** — OwnerList with toggle, OwnerProfile with career stats + history
6. **Belt Page** — Full leaderboard with bar viz, fun facts, timeline placeholder
7. **Awards Page** — Per-award history, leaderboard, superlatives
8. **Docker** — Dockerfile (multi-stage: Node build → Nginx), nginx.conf (SPA + API proxy), add to docker-compose

---

## API Gaps (workarounds for now, future endpoints recommended)

| Gap | Workaround | Future Endpoint |
|-----|-----------|-----------------|
| No belt full history | Show leaderboard only, timeline placeholder | `GET /belt/history` |
| No per-owner awards | Call all 11 season award endpoints, filter client-side | `GET /owners/:id/awards` |
| No per-owner season history | Call all 11 season standings endpoints, filter client-side | `GET /owners/:id/seasons` |
| No cross-season award summary | Call all 11 season award endpoints | `GET /awards` or `GET /awards/:type` |

---

## Verification

1. `npm run build` — compiles without errors
2. `npm run dev` — dev server starts, all pages render
3. With API running (`docker compose up`): all data loads on every page, no console errors
4. Mobile responsive: check all pages at 375px width
5. Dark theme renders correctly with all color tokens applied

