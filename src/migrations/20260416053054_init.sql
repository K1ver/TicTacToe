-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field BYTEA NOT NULL,
    status INTEGER NOT NULL,
    player_x UUID NOT NULL,
    player_o UUID NOT NULL,
    current_turn UUID NOT NULL,
    winner_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_games_player_x ON games(player_x);
CREATE INDEX idx_games_player_o ON games(player_o);
CREATE INDEX idx_games_status ON games(status);

-- Таблица leaderboard
CREATE TABLE leaderboard (
    user_id UUID PRIMARY KEY,
    login TEXT NOT NULL,
    winrate DOUBLE PRECISION NOT NULL DEFAULT 0
);

-- Инициализация leaderboard (заполняем всех пользователей)
INSERT INTO leaderboard (user_id, login)
SELECT id, login FROM users;

-- +goose Down
DROP TABLE IF EXISTS leaderboard;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
