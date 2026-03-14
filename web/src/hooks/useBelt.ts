import { useQuery } from '@tanstack/react-query';
import { getBeltCurrent, getBeltLeaderboard } from '../api/belt';

export function useBeltCurrent() {
  return useQuery({
    queryKey: ['belt', 'current'],
    queryFn: getBeltCurrent,
  });
}

export function useBeltLeaderboard() {
  return useQuery({
    queryKey: ['belt', 'leaderboard'],
    queryFn: getBeltLeaderboard,
  });
}
