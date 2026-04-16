package entity

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus int

const (
	StatusEntityIsWaiting GameStatus = iota
	StatusEntityInProgress
	StatusEntityFinished
)

type Game struct {
	ID          uuid.UUID  `db:"id"`
	Field       GameField  `db:"field"`
	Status      GameStatus `db:"status"`
	PlayerX     uuid.UUID  `db:"player_x"`
	PlayerO     uuid.UUID  `db:"player_o"`
	CurrentTurn uuid.UUID  `db:"current_turn"`
	WinnerId    uuid.UUID  `db:"winner_id"`
	CreatedAt   time.Time  `db:"created_at"`
}

type GameField []byte
