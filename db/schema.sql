ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 0;
ALTER SYSTEM SET log_error_verbosity = 'verbose';
SELECT pg_reload_conf();

-- 1. Enums
CREATE TYPE hand AS ENUM ('none','rock', 'paper', 'scissors');

-- 2. Players Table
CREATE TABLE players (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE games (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    total_rounds INTEGER NOT NULL DEFAULT 3,
    current_round INTEGER,
    player_one_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    player_two_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
    player_one_score INTEGER,
    player_two_score INTEGER,
    winner INTEGER REFERENCES players(id),
    finished BOOLEAN DEFAULT False,
    created_at timestamptz DEFAULT NOW()
);

-- CREATE TABLE player_round_input (
--     id INTEGER PRIMARY KEY,
--     player INTEGER REFERENCES players(id),
--     round_id INTEGER REFERENCES rounds(id),
--     hand_played hand NOT NUll
-- );

-- CREATE TABLE player_score (
--     id INTEGER PRIMARY KEY,
--     player INTEGER REFERENCES players(id),
--     score INTEGER NOT NULL DEFAULT 0
-- );

CREATE TABLE rounds (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    game INTEGER REFERENCES games(id) ON DELETE CASCADE,
    count INTEGER NOT NULL DEFAULT 1,
    player_one_id INTEGER REFERENCES players(id),
    player_two_id INTEGER REFERENCES players(id),
    player_one_hand hand,
    player_two_hand hand,
    winner INTEGER REFERENCES players(id),
    finished BOOLEAN DEFAULT False
);


