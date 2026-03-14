import { useParams } from 'react-router-dom';

export default function SeasonDetailPage() {
  const { year } = useParams<{ year: string }>();
  return (
    <div>
      <h1>{year} Season</h1>
      <p className="muted">Season detail — coming soon.</p>
    </div>
  );
}
