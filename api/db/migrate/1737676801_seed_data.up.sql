-- Seed data migration

-- Award Types
INSERT INTO award_types (name, description) VALUES
    ('League Champion', 'Winner of the championship game'),
    ('First Place Loser', 'Lost in the championship game'),
    ('Unlucky Award', 'Hardest schedule - highest points against'),
    ('Juggernaut Award', 'Most points scored in regular season'),
    ('Cupcake Award', 'Easiest schedule - lowest points against'),
    ('Bottomfeeder Award', 'Last place finisher');

-- Seasons
INSERT INTO seasons (year) VALUES
    (2015), (2016), (2017), (2018), (2019),
    (2020), (2021), (2022), (2023), (2024), (2025);

-- Owners
INSERT INTO owners (name, join_year, leave_year) VALUES
    ('Matt Johnson', 2015, NULL),
    ('Joe Bogdanich', 2015, NULL),
    ('Mike Sliwinski', 2016, NULL),
    ('Beth & Nick Oudin', 2015, NULL),
    ('Ron Rutkowski', 2016, NULL),
    ('Jason Cuevas', 2015, NULL),
    ('Blake Papineau', 2017, NULL),
    ('Jim Bogdanich', 2016, NULL),
    ('Ruben Rodriguez', 2016, NULL),
    ('Dan Antol', 2016, 2024),
    ('Andrew Lash', 2020, NULL),
    ('Mike Johnson', 2019, NULL),
    ('Michael Gotsch', 2015, 2019),
    ('Michael Cesario', 2015, 2018),
    ('Seth Culbreth', 2015, 2015),
    ('Jimmy Johnson', 2025, NULL),
    ('Mitch Armentrout', 2017, 2017),
    ('Jeff Cuevas', 2015, 2015),
    ('Dave Paco Vinson', 2015, 2015),
    ('Jake Hunhoff', 2016, 2016),
    ('Samantha Vinson', 2015, 2015);
