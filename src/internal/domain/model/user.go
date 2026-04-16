package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Login    string
	Password string
}

func NewUser(login, password string) *User {
	return &User{
		Id:       uuid.New(),
		Login:    login,
		Password: password,
	}
}
