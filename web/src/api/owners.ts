import { fetchApi } from './client';
import type { Owner, OwnerCareerStats } from '../types';

export function getOwners(active?: boolean): Promise<Owner[]> {
  const params = active !== undefined ? `?active=${active}` : '';
  return fetchApi<Owner[]>(`/owners${params}`);
}

export function getOwnerStats(id: number): Promise<OwnerCareerStats> {
  return fetchApi<OwnerCareerStats>(`/owners/${id}/stats`);
}
