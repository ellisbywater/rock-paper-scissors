-- 1. Enums
CREATE TABLE hand_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
);

INSERT INTO hand_types(name) VALUES ('rock'), ('paper'), ('scissors');

-- 2. Players Table
CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
);

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    total_rounds INTEGER NOT NULL DEFAULT 3,
    player_one REFERENCES NOT NULL players(id),
    player_two REFERENCES NOT NULL players(id),
    winner REFERENCES players(id),
    created_at TIMESTAMPZ DEFAULT NOW(),
);

CREATE TABLE player_round_input (
    id INTEGER PRIMARY KEY,
    player REFERENCES NOT NULL players(id),
    hand_played hand NOT NUll,
);

CREATE TABLE player_score (
    id INTEGER PRIMARY KEY,
    player REFERENCES NOT NULL players(id),
    score INTEGER NOT NULL DEFAULT 0,
);
CREATE TABLE rounds (
    id INTEGER PRIMARY KEY,
    game REFERENCES NOT NULL games(id),
    count INTEGER NOT NULL DEFAULT 1,
    player_one REFERENCES NOT NULL players(id),
    player_two REFERENCES NOT NULL players(id),
    player_one_input REFERENCES player_round_input(id),
    player_two_input REFERENCES player_round_input(id),
    winner REFERENCES players(id),
);