package model

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus int

const (
	StatusIsWaiting GameStatus = iota
	StatusInProgress
	StatusFinished
)

type Game struct {
	Id          uuid.UUID
	Field       GameField
	Status      GameStatus
	PlayerX     uuid.UUID
	PlayerO     uuid.UUID
	CurrentTurn uuid.UUID
	WinnerId    uuid.UUID
	CreatedAt   time.Time
}

func NewGame(rows, cols int) Game {
	return Game{
		Id:          uuid.New(),
		Field:       NewGameField(3, 3),
		Status:      StatusIsWaiting,
		PlayerX:     uuid.Nil,
		PlayerO:     uuid.Nil,
		CurrentTurn: uuid.Nil,
		WinnerId:    uuid.Nil,
		CreatedAt:   time.Now().Truncate(time.Second),
	}
}
