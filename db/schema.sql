

-- 1. Enums
CREATE TYPE hand AS ENUM ('rock', 'paper', 'scissors');

-- 2. Players Table
CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    total_rounds INTEGER NOT NULL DEFAULT 3,
    player_one INTEGER REFERENCES players(id),
    player_two INTEGER REFERENCES players(id),
    winner INTEGER REFERENCES players(id),
    created_at timestamptz DEFAULT NOW()
);

CREATE TABLE player_round_input (
    id INTEGER PRIMARY KEY,
    player INTEGER REFERENCES players(id),
    hand_played hand NOT NUll
);

CREATE TABLE player_score (
    id INTEGER PRIMARY KEY,
    player INTEGER REFERENCES players(id),
    score INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE rounds (
    id INTEGER PRIMARY KEY,
    game INTEGER REFERENCES games(id),
    count INTEGER NOT NULL DEFAULT 1,
    player_one INTEGER REFERENCES players(id),
    player_two INTEGER REFERENCES players(id),
    player_one_input INTEGER REFERENCES player_round_input(id),
    player_two_input INTEGER REFERENCES player_round_input(id),
    winner INTEGER REFERENCES players(id)
);