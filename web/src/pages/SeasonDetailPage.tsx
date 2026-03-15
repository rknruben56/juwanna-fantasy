import { useParams, Link } from 'react-router-dom';
import { useSeasonStandings, useSeasonAwards } from '../hooks/useSeasons';
import { StandingsTable } from '../components/seasons/StandingsTable';
import { AwardGrid } from '../components/seasons/AwardGrid';
import { Badge } from '../components/shared/Badge';
import { LoadingSpinner } from '../components/shared/LoadingSpinner';
import { ErrorState } from '../components/shared/ErrorState';
import './SeasonDetailPage.css';

export default function SeasonDetailPage() {
  const { year: yearParam } = useParams<{ year: string }>();
  const year = Number(yearParam);

  const {
    data: standings,
    isLoading: standingsLoading,
    error: standingsError,
    refetch: refetchStandings,
  } = useSeasonStandings(year);

  const {
    data: awards,
    isLoading: awardsLoading,
    error: awardsError,
    refetch: refetchAwards,
  } = useSeasonAwards(year);

  if (standingsLoading || awardsLoading) {
    return <LoadingSpinner message={`Loading ${year} season...`} />;
  }

  if (standingsError || awardsError) {
    return (
      <ErrorState
        message={`Failed to load ${year} season data.`}
        onRetry={() => {
          refetchStandings();
          refetchAwards();
        }}
      />
    );
  }

  const champion = awards?.find(
    (a) => a.award_name === 'League Champion'
  )?.owner_name;

  return (
    <div className="season-detail">
      <Link to="/seasons" className="season-detail__back">
        &larr; All Seasons
      </Link>

      <div className="season-detail__header">
        <h1>{year} Season</h1>
        {champion && (
          <div className="season-detail__champion">
            <span className="season-detail__crown">👑</span>
            <Badge text={champion} variant="gold" />
          </div>
        )}
      </div>

      {standings && standings.length > 0 && (
        <section className="season-detail__section">
          <h2>Standings</h2>
          <StandingsTable standings={standings} />
        </section>
      )}

      {awards && awards.length > 0 && (
        <section className="season-detail__section">
          <h2>Awards</h2>
          <AwardGrid awards={awards} />
        </section>
      )}
    </div>
  );
}
