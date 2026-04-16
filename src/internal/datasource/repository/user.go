package repository

import (
	"app/internal/database"
	"app/internal/datasource/entity"
	"app/internal/datasource/entity_mapper"
	"app/internal/domain/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create создание пользователя
	Create(ctx context.Context, user *entity.User) error
	// FindByLoginAndPassword поиск по логину и паролю
	FindByLoginAndPassword(ctx context.Context, login, password string) (*model.User, error)
	// GetUserById поиск пользователя по uuid
	GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error)
}

type userRepository struct {
	db *database.DataBase
}

func NewUserRepository(db *database.DataBase) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	exist, _ := u.FindByLoginAndPassword(ctx, user.Login, user.Password)
	if exist != nil {
		return fmt.Errorf("user already exist")
	}

	query := "INSERT INTO users(id, login, password) VALUES ($1, $2, $3)"

	_, err := u.db.DB.Exec(ctx, query, user.Id, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	return nil
}

func (u *userRepository) FindByLoginAndPassword(ctx context.Context, login, password string) (*model.User, error) {
	query := "SELECT id, login, password FROM users WHERE login = $1 AND password = $2"
	var user entity.User
	err := u.db.DB.QueryRow(ctx, query, login, password).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by login and password: %w", err)
	}

	return entity_mapper.UserToDomainFromEntity(&user), nil
}

func (u *userRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	query := "SELECT id, login FROM users WHERE id = $1"
	var user entity.User
	err := u.db.DB.QueryRow(ctx, query, userId).Scan(
		&user.Id,
		&user.Login,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by uuid: %w", err)
	}

	return entity_mapper.UserToDomainFromEntity(&user), nil
}
