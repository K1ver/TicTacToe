package auth

import (
	"app/internal/domain/service"
	"context"
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	jwt         service.JwtProvider
	userService service.UserService
}

func NewUserAuthenticator(jwt service.JwtProvider, userService service.UserService) UserAuthenticator {
	return UserAuthenticator{
		jwt:         jwt,
		userService: userService,
	}
}

func (u *UserAuthenticator) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "bearer authentication required", http.StatusUnauthorized)
				return
			}
			encoded := strings.TrimPrefix(authHeader, "Bearer ")
			err := u.jwt.ValidateAccessToken(context.Background(), encoded)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userId, err := u.jwt.GetUserIdByToken(context.Background(), encoded)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			_, err = u.userService.GetUserById(context.Background(), userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Добавляем userID в контекст запроса
			ctx := context.WithValue(r.Context(), "userID", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
