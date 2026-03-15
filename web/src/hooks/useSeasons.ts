import { useQuery } from '@tanstack/react-query';
import { getSeasonStandings, getSeasonAwards } from '../api/seasons';

export function useSeasonStandings(year: number) {
  return useQuery({
    queryKey: ['seasons', year, 'standings'],
    queryFn: () => getSeasonStandings(year),
  });
}

export function useSeasonAwards(year: number) {
  return useQuery({
    queryKey: ['seasons', year, 'awards'],
    queryFn: () => getSeasonAwards(year),
  });
}
