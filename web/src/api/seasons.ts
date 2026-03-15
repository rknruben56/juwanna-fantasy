import { fetchApi } from './client';
import type { SeasonStanding, SeasonAwardWithDetails } from '../types';

export function getSeasonStandings(year: number): Promise<SeasonStanding[]> {
  return fetchApi<SeasonStanding[]>(`/seasons/${year}/standings`);
}

export function getSeasonAwards(year: number): Promise<SeasonAwardWithDetails[]> {
  return fetchApi<SeasonAwardWithDetails[]>(`/seasons/${year}/awards`);
}
