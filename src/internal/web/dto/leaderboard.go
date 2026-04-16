package dto

import "github.com/google/uuid"

type Leaderboard struct {
	UserId  uuid.UUID `json:"userId"`
	Login   string    `json:"login"`
	Winrate float64   `json:"winrate"`
}
