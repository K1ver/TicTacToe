package service

import (
	"app/internal/datasource/entity_mapper"
	"app/internal/datasource/repository"
	"app/internal/domain/model"
	"app/internal/web/dto"
	"app/internal/web/dto_mapper"
	"context"
	"log"

	"github.com/google/uuid"
)

type UserService interface {
	// SignUp метод регистрации
	SignUp(ctx context.Context, request *model.JwtRequest) (bool, error)
	// SignIn метод авторизации
	SignIn(ctx context.Context, request *model.JwtRequest) (*dto.JwtResponse, error)
	// GetUserById поиск пользователя по uuid
	GetUserById(ctx context.Context, userId uuid.UUID) (*dto.User, error)
	// GetUserByToken поиск пользователя по accessToken
	GetUserByToken(ctx context.Context, tokenString string) (*dto.User, error)
	// UpdateAccessToken обновление accessToken
	UpdateAccessToken(ctx context.Context, request *model.RefreshJwtRequest) (*dto.JwtResponse, error)
	// UpdateRefreshToken обновление refreshToken
	UpdateRefreshToken(ctx context.Context, request *model.RefreshJwtRequest) (*dto.JwtResponse, error)
}

type userService struct {
	repo repository.UserRepository
	jwt  JwtProvider
}

func NewUserService(repo repository.UserRepository, jwt JwtProvider) UserService {
	return &userService{
		repo: repo,
		jwt:  jwt,
	}
}

func (u *userService) SignUp(ctx context.Context, request *model.JwtRequest) (bool, error) {
	newUser := model.NewUser(request.Login, request.Password)
	err := u.repo.Create(ctx, entity_mapper.UserToEntity(newUser))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func (u *userService) SignIn(ctx context.Context, request *model.JwtRequest) (*dto.JwtResponse, error) {
	login, password := request.Login, request.Password
	user, err := u.repo.FindByLoginAndPassword(ctx, login, password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	accessToken, err := u.jwt.GenerateAccessToken(ctx, user.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateRefreshToken(ctx, user.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (u *userService) GetUserById(ctx context.Context, userId uuid.UUID) (*dto.User, error) {
	user, err := u.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return dto_mapper.UserToDto(user), nil
}

func (u *userService) GetUserByToken(ctx context.Context, tokenString string) (*dto.User, error) {
	userId, err := u.jwt.GetUserIdByToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	return u.GetUserById(ctx, userId)
}

func (u *userService) UpdateAccessToken(ctx context.Context, request *model.RefreshJwtRequest) (*dto.JwtResponse, error) {
	err := u.jwt.ValidateRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	userId, err := u.jwt.GetUserIdByToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwt.GenerateAccessToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: request.RefreshToken,
	}, nil
}
func (u *userService) UpdateRefreshToken(ctx context.Context, request *model.RefreshJwtRequest) (*dto.JwtResponse, error) {
	err := u.jwt.ValidateRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	userId, err := u.jwt.GetUserIdByToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwt.GenerateRefreshToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.JwtResponse{
		AccessToken:  "",
		RefreshToken: refreshToken,
	}, nil
}
