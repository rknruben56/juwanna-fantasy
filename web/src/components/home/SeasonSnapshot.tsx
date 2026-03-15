import { Link } from 'react-router-dom';
import { useSeasonStandings, useSeasonAwards } from '../../hooks/useSeasons';
import { DataTable, Badge, AwardIcon, LoadingSpinner, ErrorState } from '../shared';
import type { Column } from '../shared';
import type { SeasonStanding } from '../../types';
import './SeasonSnapshot.css';

const LATEST_SEASON = 2025;

const columns: Column<SeasonStanding & Record<string, unknown>>[] = [
  {
    key: 'rank',
    label: '#',
    align: 'center',
    render: (_val, _row, ) => null, // overridden below
  },
  { key: 'owner_name', label: 'Owner', sortable: true },
  {
    key: 'regular_season_wins',
    label: 'W-L',
    align: 'center',
    render: (_val, row) => `${row.regular_season_wins}-${row.regular_season_losses}`,
  },
  {
    key: 'made_playoffs',
    label: 'Playoffs',
    align: 'center',
    render: (val) => val ? <Badge text="IN" variant="teal" size="sm" /> : null,
  },
];

export function SeasonSnapshot() {
  const { data: standings, isLoading: loadingStandings, error: standingsError, refetch: refetchStandings } = useSeasonStandings(LATEST_SEASON);
  const { data: awards, isLoading: loadingAwards } = useSeasonAwards(LATEST_SEASON);

  if (loadingStandings || loadingAwards) return <LoadingSpinner message="Loading latest season..." />;
  if (standingsError) return <ErrorState message="Failed to load season data" onRetry={refetchStandings} />;
  if (!standings?.length) return null;

  // Sort by wins desc for ranking
  const sorted = [...standings].sort((a, b) => b.regular_season_wins - a.regular_season_wins);

  // Add rank via render override
  const rankedColumns = columns.map((col) =>
    col.key === 'rank'
      ? { ...col, render: (_val: unknown, row: SeasonStanding & Record<string, unknown>) => sorted.indexOf(row as unknown as SeasonStanding & Record<string, unknown>) + 1 }
      : col
  );

  return (
    <div className="season-snapshot">
      <DataTable
        columns={rankedColumns}
        data={sorted as (SeasonStanding & Record<string, unknown>)[]}
        defaultSort={{ key: 'regular_season_wins', direction: 'desc' }}
        rowClassName={(row) =>
          row.made_playoffs ? 'season-snapshot__row--playoffs' : ''
        }
      />

      {awards && awards.length > 0 && (
        <div className="season-snapshot__awards">
          {awards.map((award) => (
            <div key={award.id} className="season-snapshot__award">
              <AwardIcon awardName={award.award_name} size="md" />
              <div className="season-snapshot__award-info">
                <span className="season-snapshot__award-name">{award.award_name}</span>
                <span className="season-snapshot__award-owner">{award.owner_name}</span>
              </div>
            </div>
          ))}
        </div>
      )}

      <Link to={`/seasons/${LATEST_SEASON}`} className="season-snapshot__link">
        View Full Season →
      </Link>
    </div>
  );
}
