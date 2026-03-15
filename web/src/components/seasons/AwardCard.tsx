import { AwardIcon } from '../shared/AwardIcon';
import type { SeasonAwardWithDetails } from '../../types';
import './AwardCard.css';

interface AwardCardProps {
  award: SeasonAwardWithDetails;
}

export function AwardCard({ award }: AwardCardProps) {
  const isChampion = award.award_name === 'League Champion';
  const isBottomfeeder = award.award_name === 'Bottomfeeder';

  const cardClass = [
    'award-card',
    isChampion ? 'award-card--champion' : '',
    isBottomfeeder ? 'award-card--bottomfeeder' : '',
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={cardClass}>
      <AwardIcon awardName={award.award_name} size="lg" />
      <div className="award-card__info">
        <div className="award-card__name">{award.award_name}</div>
        <div className="award-card__owner">{award.owner_name}</div>
        {award.stat_value != null && (
          <div className="award-card__stat">
            {Number(award.stat_value).toFixed(1)} pts
          </div>
        )}
      </div>
    </div>
  );
}
