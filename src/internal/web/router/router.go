package router

import (
	"app/internal/web/auth"
	"app/internal/web/handler"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

func RegisterRoutes(lc fx.Lifecycle, gameHandler handler.GameHandler, userHandler handler.UserHandler, authenticator auth.UserAuthenticator) {
	r := chi.NewRouter()

	// Базовые middleware - ДО любых маршрутов
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Сначала API маршруты
	r.Route("/api", func(r chi.Router) {

		// Маршруты авторизации (публичные)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signUp", userHandler.SignUpHandle)
			r.Post("/signIn", userHandler.SignInHandle)
		})

		// Маршруты пользователя
		r.Route("/user", func(r chi.Router) {
			// Требуют авторизации
			r.Use(authenticator.AuthMiddleware())
			r.Get("/me", userHandler.GetUserByTokenHandle)
			r.Get("/{userId}", userHandler.GetUserByIdHandle)
		})

		// Маршруты для токена
		r.Route("/token", func(r chi.Router) {
			// Обновление access токена - публичный
			r.Post("/updateAccess", userHandler.UpdateAccessTokenHandle)

			// Защищенные маршруты
			r.Group(func(r chi.Router) {
				r.Use(authenticator.AuthMiddleware())
				r.Post("/updateRefresh", userHandler.UpdateRefreshTokenHandle)
			})
		})

		// Маршруты игры
		r.Route("/game", func(r chi.Router) {
			//Требуют авторизации
			r.Use(authenticator.AuthMiddleware())

			r.Post("/", gameHandler.CreateGameHandle)
			r.Post("/{gameId}/join", gameHandler.JoinGameHandle)
			r.Post("/{gameId}", gameHandler.MakeMoveHandle)
			r.Get("/{gameId}", gameHandler.GetGameHandle)
			r.Get("/history", gameHandler.GetAllFinishedGamesHandle)
			r.Get("/browse", gameHandler.GetAllUnstartedGamesHandle)
			r.Get("/leaderboard/{count}", gameHandler.GetLeaderboardHandle)
		})
	})

	// Статические файлы фронтенда
	r.Handle("/*", http.FileServer(http.Dir("./web/")))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				_ = server.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
