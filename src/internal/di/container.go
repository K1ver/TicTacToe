package di

import (
	"app/internal/database"
	"app/internal/datasource/repository"
	"app/internal/domain/algorithm"
	"app/internal/domain/service"
	"app/internal/web/auth"
	"app/internal/web/handler"
	"app/internal/web/router"
	"context"

	"go.uber.org/fx"
)

// Module модуль зависимостей
var Module = fx.Module("app",
	// Инфраструктура
	fx.Provide(
		database.NewDataBase,
		repository.NewGameRepository,
		repository.NewUserRepository,
	),

	// Domain
	fx.Provide(
		fx.Annotate(
			algorithm.NewMinimax,
			fx.As(new(algorithm.MinimaxAlgorithm)),
		),
	),

	// Application
	fx.Provide(
		fx.Annotate(
			service.NewJwtProvider,
			fx.As(new(service.JwtProvider)),
		),
		fx.Annotate(
			service.NewGameService,
			fx.As(new(service.GameService)),
		),
		fx.Annotate(
			service.NewUserService,
			fx.As(new(service.UserService)),
		),
	),

	// Web
	fx.Provide(
		fx.Annotate(
			handler.NewGameHandler,
			fx.As(new(handler.GameHandler)),
		),
		fx.Annotate(
			handler.NewUserHandler,
			fx.As(new(handler.UserHandler)),
		),
		auth.NewUserAuthenticator,
	),

	//Маршрутизатор
	fx.Invoke(router.RegisterRoutes),

	// Жизненный цикл
	fx.Invoke(registerHooks),
)

func registerHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Инициализация при запуске
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Очистка при остановке
			return nil
		},
	})
}
