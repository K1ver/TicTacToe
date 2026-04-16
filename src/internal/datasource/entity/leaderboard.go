package entity

import "github.com/google/uuid"

type Leaderboard struct {
	UserId  uuid.UUID `db:"user_id"`
	Login   string    `db:"login"`
	Winrate float64   `db:"winrate"`
}
