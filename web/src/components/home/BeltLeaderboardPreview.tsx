import { Link } from 'react-router-dom';
import { useBeltLeaderboard } from '../../hooks/useBelt';
import { LoadingSpinner, ErrorState } from '../shared';
import './BeltLeaderboardPreview.css';

export function BeltLeaderboardPreview() {
  const { data: leaderboard, isLoading, error, refetch } = useBeltLeaderboard();

  if (isLoading) return <LoadingSpinner message="Loading belt leaderboard..." />;
  if (error) return <ErrorState message="Failed to load leaderboard" onRetry={refetch} />;
  if (!leaderboard?.length) return null;

  const top5 = leaderboard.slice(0, 5);
  const maxWeeks = top5[0].weeks_with_belt;

  return (
    <div className="belt-preview">
      <div className="belt-preview__list">
        {top5.map((entry, i) => (
          <div key={entry.owner_id} className="belt-preview__row">
            <span className="belt-preview__rank">{i + 1}</span>
            <span className="belt-preview__name">{entry.name}</span>
            <div className="belt-preview__bar-track">
              <div
                className="belt-preview__bar-fill"
                style={{ width: `${(entry.weeks_with_belt / maxWeeks) * 100}%` }}
              />
            </div>
            <span className="belt-preview__weeks">{entry.weeks_with_belt}w</span>
          </div>
        ))}
      </div>
      <Link to="/belt" className="belt-preview__link">View Full Leaderboard →</Link>
    </div>
  );
}
