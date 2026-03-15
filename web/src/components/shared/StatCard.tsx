import type { ReactNode } from 'react';
import './StatCard.css';

interface StatCardProps {
  label: string;
  value: string | number;
  color?: 'gold' | 'teal' | 'red' | 'default';
  icon?: ReactNode;
}

export function StatCard({ label, value, color = 'default', icon }: StatCardProps) {
  return (
    <div className="stat-card">
      {icon && <div className="stat-card__icon">{icon}</div>}
      <div className={`stat-card__value stat-card__value--${color}`}>{value}</div>
      <div className="stat-card__label">{label}</div>
    </div>
  );
}
