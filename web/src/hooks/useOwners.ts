import { useQuery } from '@tanstack/react-query';
import { getOwners, getOwnerStats } from '../api/owners';

export function useOwners(active?: boolean) {
  return useQuery({
    queryKey: ['owners', { active }],
    queryFn: () => getOwners(active),
  });
}

export function useOwnerStats(id: number) {
  return useQuery({
    queryKey: ['owners', id, 'stats'],
    queryFn: () => getOwnerStats(id),
  });
}
