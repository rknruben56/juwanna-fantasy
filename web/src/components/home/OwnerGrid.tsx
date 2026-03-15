import { Link } from 'react-router-dom';
import { useOwners } from '../../hooks/useOwners';
import { LoadingSpinner, ErrorState } from '../shared';
import './OwnerGrid.css';

export function OwnerGrid() {
  const { data: owners, isLoading, error, refetch } = useOwners(true);

  if (isLoading) return <LoadingSpinner message="Loading owners..." />;
  if (error) return <ErrorState message="Failed to load owners" onRetry={refetch} />;
  if (!owners?.length) return null;

  return (
    <div className="owner-grid">
      {owners.map((owner) => (
        <Link key={owner.id} to={`/owners/${owner.id}`} className="owner-grid__card">
          <div className="owner-grid__avatar">
            {owner.name.charAt(0).toUpperCase()}
          </div>
          <h3 className="owner-grid__name">{owner.name}</h3>
          <p className="owner-grid__since">Since {owner.join_year}</p>
        </Link>
      ))}
    </div>
  );
}
