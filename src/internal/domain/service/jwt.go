package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type JwtProvider interface {
	// GenerateAccessToken генерация AccessToken по id пользователя
	GenerateAccessToken(ctx context.Context, userId uuid.UUID) (string, error)
	// GenerateRefreshToken генерация RefreshToken по id пользователя
	GenerateRefreshToken(ctx context.Context, userId uuid.UUID) (string, error)
	// ValidateAccessToken проверка подлинности AccessToken
	ValidateAccessToken(ctx context.Context, accessToken string) error
	// ValidateRefreshToken проверка подлинности RefreshToken
	ValidateRefreshToken(ctx context.Context, refreshToken string) error
	// GetUserIdByToken извлечение userId из токена
	GetUserIdByToken(ctx context.Context, accessToken string) (uuid.UUID, error)
}

type jwtProvider struct {
	secretKey []byte
}

func NewJwtProvider() JwtProvider {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	return &jwtProvider{
		secretKey: secretKey,
	}
}

func (j *jwtProvider) GenerateAccessToken(ctx context.Context, userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 20).Unix(), // Срок действия 20 минут
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)

}

func (j *jwtProvider) GenerateRefreshToken(ctx context.Context, userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 50).Unix(), // Срок действия 50 минут
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtProvider) ValidateAccessToken(ctx context.Context, accessToken string) error {
	return j.validateToken(j.parseToken(accessToken))
}

func (j *jwtProvider) ValidateRefreshToken(ctx context.Context, refreshToken string) error {
	return j.validateToken(j.parseToken(refreshToken))
}

func (j *jwtProvider) GetUserIdByToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	token, err := j.parseToken(accessToken)
	err = j.validateToken(token, err)
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		str := claims["userId"]
		return uuid.Parse(str.(string))
	} else {
		return uuid.Nil, fmt.Errorf("hz")
	}
}

func (j *jwtProvider) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.secretKey, nil
	})
}

func (j *jwtProvider) validateToken(token *jwt.Token, err error) error {
	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf("that's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// Invalid signature
		return fmt.Errorf("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// Token is either expired or not active yet
		return fmt.Errorf("timing is everything")
	default:
		return fmt.Errorf("couldn't handle this token" + err.Error())
	}
}
