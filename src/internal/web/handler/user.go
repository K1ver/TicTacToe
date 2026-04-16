package handler

import (
	"app/internal/domain/service"
	"app/internal/web/dto"
	"app/internal/web/dto_mapper"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler interface {
	// SignUpHandle хендлер регистрации
	SignUpHandle(cw http.ResponseWriter, r *http.Request)
	// SignInHandle хендлер авторизации
	SignInHandle(w http.ResponseWriter, r *http.Request)
	// GetUserByIdHandle получение информации о пользователе по его UUID
	GetUserByIdHandle(w http.ResponseWriter, r *http.Request)
	// GetUserByTokenHandle получение информации о пользователе по accessToken
	GetUserByTokenHandle(w http.ResponseWriter, r *http.Request)
	// UpdateAccessTokenHandle обновление accessToken
	UpdateAccessTokenHandle(w http.ResponseWriter, r *http.Request)
	// UpdateRefreshTokenHandle обновление refreshToken
	UpdateRefreshTokenHandle(w http.ResponseWriter, r *http.Request)
}

type userHandle struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandle{userService: userService}
}

func (u *userHandle) SignUpHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	var req dto.JwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := u.userService.SignUp(ctx, dto_mapper.JwtRequestToDomainFromDto(&req))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *userHandle) SignInHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	var req dto.JwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := u.userService.SignIn(ctx, dto_mapper.JwtRequestToDomainFromDto(&req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *userHandle) GetUserByIdHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	userIdStr := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := u.userService.GetUserById(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *userHandle) GetUserByTokenHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "authorization header required", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "bearer authentication required", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := u.userService.GetUserByToken(ctx, tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

}

func (u *userHandle) UpdateAccessTokenHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var req dto.RefreshJwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := u.userService.UpdateAccessToken(ctx, dto_mapper.RefreshJwtRequestToDomainFromDto(&req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (u *userHandle) UpdateRefreshTokenHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var req dto.RefreshJwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := u.userService.UpdateRefreshToken(ctx, dto_mapper.RefreshJwtRequestToDomainFromDto(&req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
