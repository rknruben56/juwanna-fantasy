import { Link } from 'react-router-dom';
import type { SeasonStanding } from '../../types';
import './SeasonCard.css';

interface SeasonCardProps {
  year: number;
  champion?: string;
  standings?: SeasonStanding[];
}

export function SeasonCard({ year, champion, standings }: SeasonCardProps) {
  const top3 = standings?.slice(0, 3);

  return (
    <Link to={`/seasons/${year}`} className="season-card">
      <div className="season-card__year">{year}</div>
      {champion && (
        <div className="season-card__champion">
          <span className="season-card__crown">👑</span>
          <span className="season-card__champion-name">{champion}</span>
        </div>
      )}
      {top3 && top3.length > 0 && (
        <ol className="season-card__standings">
          {top3.map((s) => (
            <li key={s.owner_name} className="season-card__standing-row">
              <span className="season-card__owner">{s.owner_name}</span>
              <span className="season-card__record">
                {s.regular_season_wins}-{s.regular_season_losses}
              </span>
            </li>
          ))}
        </ol>
      )}
    </Link>
  );
}
