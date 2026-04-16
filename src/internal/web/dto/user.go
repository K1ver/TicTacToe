package dto

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}
