import { Link } from 'react-router-dom';

export default function NotFoundPage() {
  return (
    <div>
      <h1>404</h1>
      <p className="muted">Page not found.</p>
      <p><Link to="/">Back to home</Link></p>
    </div>
  );
}
