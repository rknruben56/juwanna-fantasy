// Core domain models

export interface Owner {
  id: number;
  name: string;
  join_year: number;
  leave_year?: number;
  created_at: string;
  updated_at: string;
}

export interface OwnerCareerStats {
  owner_id: number;
  name: string;
  join_year: number;
  leave_year?: number;
  years_in_league: number;
  total_regular_season_wins: number;
  total_regular_season_losses: number;
  total_playoff_wins: number;
  total_playoff_losses: number;
  playoff_appearances: number;
  championships: number;
  first_place_loser_trophies: number;
}

export interface OwnerBeltWeeks {
  owner_id: number;
  name: string;
  weeks_with_belt: number;
}

export interface SeasonStanding {
  year: number;
  owner_name: string;
  regular_season_wins: number;
  regular_season_losses: number;
  points_for?: number;
  points_against?: number;
  made_playoffs: boolean;
  team_name?: string;
}

export interface SeasonAwardWithDetails {
  id: number;
  owner_name: string;
  award_name: string;
  year: number;
  stat_value?: number;
  team_name?: string;
}

export interface BeltHistoryWithDetails {
  owner_name: string;
  year: number;
  week_number: number;
}
