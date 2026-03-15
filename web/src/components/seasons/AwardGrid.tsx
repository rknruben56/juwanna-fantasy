import { AwardCard } from './AwardCard';
import type { SeasonAwardWithDetails } from '../../types';
import './AwardGrid.css';

interface AwardGridProps {
  awards: SeasonAwardWithDetails[];
}

export function AwardGrid({ awards }: AwardGridProps) {
  const sorted = [...awards].sort((a, b) => {
    if (a.award_name === 'League Champion') return -1;
    if (b.award_name === 'League Champion') return 1;
    return a.award_name.localeCompare(b.award_name);
  });

  return (
    <div className="award-grid">
      {sorted.map((award) => (
        <AwardCard key={award.id} award={award} />
      ))}
    </div>
  );
}
