import { BeltHero } from '../components/home/BeltHero';
import { QuickStatsBar } from '../components/home/QuickStatsBar';
import { BeltLeaderboardPreview } from '../components/home/BeltLeaderboardPreview';
import { SeasonSnapshot } from '../components/home/SeasonSnapshot';
import { OwnerGrid } from '../components/home/OwnerGrid';
import './HomePage.css';

export default function HomePage() {
  return (
    <div className="home">
      <BeltHero />

      <section className="home__section">
        <h2 className="home__heading">League at a Glance</h2>
        <QuickStatsBar />
      </section>

      <section className="home__section">
        <h2 className="home__heading">Belt Leaderboard</h2>
        <BeltLeaderboardPreview />
      </section>

      <section className="home__section">
        <h2 className="home__heading">Latest Season — 2025</h2>
        <SeasonSnapshot />
      </section>

      <section className="home__section">
        <h2 className="home__heading">Active Owners</h2>
        <OwnerGrid />
      </section>
    </div>
  );
}
