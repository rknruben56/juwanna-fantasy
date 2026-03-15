import { useBeltCurrent } from '../../hooks/useBelt';
import { LoadingSpinner } from '../shared';
import './BeltHero.css';

export function BeltHero() {
  const { data: belt, isLoading, error } = useBeltCurrent();

  if (isLoading) return <LoadingSpinner message="Loading belt holder..." />;
  if (error || !belt) {
    return (
      <section className="belt-hero">
        <div className="belt-hero__content">
          <div className="belt-hero__belt-icon">🏆</div>
          <p className="belt-hero__subtitle">Championship Belt</p>
          <h2 className="belt-hero__name" style={{ color: 'var(--text-muted)' }}>Unclaimed</h2>
          <p className="belt-hero__since">No belt history recorded yet</p>
        </div>
      </section>
    );
  }

  return (
    <section className="belt-hero">
      <div className="belt-hero__glow" />
      <div className="belt-hero__content">
        <div className="belt-hero__belt-icon">🏆</div>
        <p className="belt-hero__subtitle">Current Championship Belt Holder</p>
        <h2 className="belt-hero__name">{belt.owner_name}</h2>
        <p className="belt-hero__since">
          Held since Week {belt.week_number}, {belt.year}
        </p>
      </div>
    </section>
  );
}
