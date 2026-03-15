import './AwardIcon.css';

const AWARD_MAP: Record<string, { emoji: string; labelClass: string }> = {
  'League Champion': { emoji: '👑', labelClass: 'award-icon__label--champion' },
  'First Place Loser': { emoji: '🥈', labelClass: '' },
  'Unlucky': { emoji: '⛈️', labelClass: 'award-icon__label--unlucky' },
  'Juggernaut': { emoji: '🚀', labelClass: 'award-icon__label--juggernaut' },
  'Cupcake': { emoji: '🧁', labelClass: 'award-icon__label--cupcake' },
  'Bottomfeeder': { emoji: '🗑️', labelClass: 'award-icon__label--bottomfeeder' },
};

interface AwardIconProps {
  awardName: string;
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
}

export function AwardIcon({ awardName, size = 'md', showLabel = false }: AwardIconProps) {
  const award = AWARD_MAP[awardName] ?? { emoji: '🏆', labelClass: '' };

  return (
    <span className="award-icon" title={awardName}>
      <span className={`award-icon__emoji--${size}`}>{award.emoji}</span>
      {showLabel && (
        <span className={`award-icon__label ${award.labelClass}`}>{awardName}</span>
      )}
    </span>
  );
}
