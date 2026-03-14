import { NavLink } from 'react-router-dom';
import './Navbar.css';

export default function Navbar() {
  return (
    <nav className="navbar">
      <div className="container navbar-inner">
        <NavLink to="/" className="navbar-logo">
          Juwanna Fantasy
        </NavLink>
        <div className="navbar-links">
          <NavLink to="/" end>Home</NavLink>
          <NavLink to="/owners">Owners</NavLink>
          <NavLink to="/seasons">Seasons</NavLink>
          <NavLink to="/belt">Belt</NavLink>
          <NavLink to="/awards">Awards</NavLink>
        </div>
      </div>
    </nav>
  );
}
