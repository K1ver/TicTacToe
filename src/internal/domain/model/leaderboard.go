package model

import "github.com/google/uuid"

type Leaderboard struct {
	UserId  uuid.UUID
	Login   string
	Winrate float64
}
