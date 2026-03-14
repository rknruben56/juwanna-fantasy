import { useParams } from 'react-router-dom';

export default function OwnerProfilePage() {
  const { id } = useParams<{ id: string }>();
  return (
    <div>
      <h1>Owner Profile</h1>
      <p className="muted">Owner #{id} — coming soon.</p>
    </div>
  );
}
