import { fetchApi } from './client';
import type { BeltHistoryWithDetails, OwnerBeltWeeks } from '../types';

export function getBeltCurrent(): Promise<BeltHistoryWithDetails> {
  return fetchApi<BeltHistoryWithDetails>('/belt/current');
}

export function getBeltLeaderboard(): Promise<OwnerBeltWeeks[]> {
  return fetchApi<OwnerBeltWeeks[]>('/belt/leaderboard');
}
