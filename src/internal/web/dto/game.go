package dto

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus int

const (
	StatusDtoIsWaiting GameStatus = iota
	StatusDtoInProgress
	StatusDtoFinished
)

type Game struct {
	Id          uuid.UUID  `json:"id"`
	Field       GameField  `json:"field"`
	Status      GameStatus `json:"status"`
	PlayerX     uuid.UUID  `json:"playerX"`
	PlayerO     uuid.UUID  `json:"playerO"`
	CurrentTurn uuid.UUID  `json:"currentTurn"`
	WinnerId    uuid.UUID  `json:"winnerId"`
	CreatedAt   time.Time  `json:"createdAt"`
}
