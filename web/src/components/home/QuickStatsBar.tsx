import { useOwners } from '../../hooks/useOwners';
import { useBeltLeaderboard } from '../../hooks/useBelt';
import { StatCard, LoadingSpinner } from '../shared';
import './QuickStatsBar.css';

export function QuickStatsBar() {
  const { data: allOwners, isLoading: loadingAll } = useOwners();
  const { data: activeOwners, isLoading: loadingActive } = useOwners(true);
  const { data: beltLeaderboard, isLoading: loadingBelt } = useBeltLeaderboard();

  if (loadingAll || loadingActive || loadingBelt) {
    return <LoadingSpinner message="Loading stats..." />;
  }

  const totalBeltWeeks = beltLeaderboard?.reduce((sum, o) => sum + o.weeks_with_belt, 0) ?? 0;

  return (
    <div className="quick-stats-bar">
      <StatCard label="Seasons" value={11} icon={<span>📅</span>} />
      <StatCard label="All-Time Owners" value={allOwners?.length ?? 0} icon={<span>👥</span>} />
      <StatCard label="Active Owners" value={activeOwners?.length ?? 0} color="teal" icon={<span>🏈</span>} />
      <StatCard label="Total Belt Weeks" value={totalBeltWeeks} color="gold" icon={<span>🏆</span>} />
    </div>
  );
}
