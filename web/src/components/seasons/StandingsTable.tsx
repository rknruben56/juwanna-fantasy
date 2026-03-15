import { DataTable } from '../shared/DataTable';
import type { Column } from '../shared/DataTable';
import { Badge } from '../shared/Badge';
import type { SeasonStanding } from '../../types';
import './StandingsTable.css';

interface StandingsTableProps {
  standings: SeasonStanding[];
}

interface StandingRow extends Record<string, unknown> {
  rank: number;
  owner_name: string;
  team_name?: string;
  regular_season_wins: number;
  regular_season_losses: number;
  points_for?: number;
  points_against?: number;
  made_playoffs: boolean;
  diff: number | null;
  year: number;
}

const COLUMNS: Column<StandingRow>[] = [
  {
    key: 'rank',
    label: '#',
    align: 'center',
  },
  {
    key: 'owner_name',
    label: 'Owner',
    sortable: true,
  },
  {
    key: 'team_name',
    label: 'Team',
    render: (val) => (val ? String(val) : '\u2014'),
  },
  {
    key: 'regular_season_wins',
    label: 'W-L',
    sortable: true,
    align: 'center',
    render: (_val, row) =>
      `${row.regular_season_wins}-${row.regular_season_losses}`,
  },
  {
    key: 'points_for',
    label: 'PF',
    sortable: true,
    align: 'right',
    render: (val) => (val != null ? Number(val).toFixed(1) : '\u2014'),
  },
  {
    key: 'points_against',
    label: 'PA',
    sortable: true,
    align: 'right',
    render: (val) => (val != null ? Number(val).toFixed(1) : '\u2014'),
  },
  {
    key: 'diff',
    label: '+/\u2212',
    sortable: true,
    align: 'right',
    render: (val) => {
      if (val == null) return '\u2014';
      const n = Number(val);
      const sign = n > 0 ? '+' : '';
      return `${sign}${n.toFixed(1)}`;
    },
  },
  {
    key: 'made_playoffs',
    label: 'Playoffs',
    align: 'center',
    render: (val) =>
      val ? (
        <Badge text="Playoffs" variant="teal" size="sm" />
      ) : (
        <Badge text="Missed" variant="muted" size="sm" />
      ),
  },
];

export function StandingsTable({ standings }: StandingsTableProps) {
  const rows: StandingRow[] = standings.map((s, i) => ({
    ...s,
    rank: i + 1,
    diff:
      s.points_for != null && s.points_against != null
        ? s.points_for - s.points_against
        : null,
  }));

  const lastRank = rows.length;

  return (
    <div className="standings-table">
      <DataTable<StandingRow>
        columns={COLUMNS}
        data={rows}
        rowClassName={(row) => {
          if (row.rank === 1) return 'standings-row--champion';
          if (row.rank === lastRank) return 'standings-row--last';
          if (row.made_playoffs) return 'standings-row--playoffs';
          return '';
        }}
      />
    </div>
  );
}
