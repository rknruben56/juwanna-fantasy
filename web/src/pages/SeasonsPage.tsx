import { useSeasonStandings, useSeasonAwards } from '../hooks/useSeasons';
import { SeasonCard } from '../components/seasons/SeasonCard';
import './SeasonsPage.css';

const YEARS = Array.from({ length: 11 }, (_, i) => 2025 - i);

function SeasonCardWithData({ year }: { year: number }) {
  const { data: standings } = useSeasonStandings(year);
  const { data: awards } = useSeasonAwards(year);

  const champion = awards?.find((a) => a.award_name === 'League Champion')?.owner_name;

  return <SeasonCard year={year} champion={champion} standings={standings} />;
}

export default function SeasonsPage() {
  return (
    <div>
      <h1>Season Archives</h1>
      <p className="muted" style={{ marginBottom: 'var(--space-2xl)' }}>
        11 seasons of glory, heartbreak, and questionable roster decisions.
      </p>
      <div className="seasons-grid">
        {YEARS.map((year) => (
          <SeasonCardWithData key={year} year={year} />
        ))}
      </div>
    </div>
  );
}
