package repository

import (
	"app/internal/database"
	"app/internal/datasource/entity"
	"app/internal/datasource/entity_mapper"
	"app/internal/domain/model"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type GameRepository interface {
	// Create создание игры
	Create(ctx context.Context, game *entity.Game) error
	// Update обновление игры
	Update(ctx context.Context, game *entity.Game) error
	// Delete удаление игры
	Delete(ctx context.Context, id uuid.UUID) error
	// GetByID получение игры по ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	// GetAllByStatus получение всех игр по статусу
	GetAllByStatus(ctx context.Context, status entity.GameStatus) ([]*model.Game, error)
	// GetFinishedByUserId получения всех завершенных игр по UUID пользователя.
	GetFinishedByUserId(ctx context.Context, userId uuid.UUID) ([]*model.Game, error)
	//GetLeaderboard получение N записей пользователей и соотношение побед и поражений
	GetLeaderboard(ctx context.Context, count int) ([]*model.Leaderboard, error)
}

type gameRepository struct {
	db *database.DataBase
}

func NewGameRepository(db *database.DataBase) GameRepository {
	return &gameRepository{
		db: db,
	}
}

func (g gameRepository) Create(ctx context.Context, game *entity.Game) error {
	log.Print(game.CreatedAt)
	query := "INSERT INTO games(id, field, status, player_x, player_o, current_turn, winner_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := g.db.DB.Exec(ctx, query, game.ID, game.Field, game.Status, game.PlayerX, game.PlayerO, game.CurrentTurn, game.WinnerId, game.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}
	return nil
}

func (g gameRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	query := "SELECT id, field, status, player_x, player_o, current_turn, winner_id, created_at FROM games WHERE id = $1"
	var game entity.Game
	var fieldJSON []byte
	err := g.db.DB.QueryRow(ctx, query, id).Scan(
		&game.ID,
		&fieldJSON,
		&game.Status,
		&game.PlayerX,
		&game.PlayerO,
		&game.CurrentTurn,
		&game.WinnerId,
		&game.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("not found game: %w", err)
	}

	game.Field = fieldJSON

	return entity_mapper.GameToDomainFromEntity(&game), nil
}

func (g gameRepository) GetAllByStatus(ctx context.Context, status entity.GameStatus) ([]*model.Game, error) {
	query := "SELECT id, field, status, player_x, player_o, current_turn, winner_id, created_at FROM games WHERE status = $1"

	rows, err := g.db.DB.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("wrong query")
	}
	var eGames []entity.Game
	for rows.Next() {
		var eGame entity.Game
		var fieldJSON []byte

		err := rows.Scan(
			&eGame.ID,
			&fieldJSON,
			&eGame.Status,
			&eGame.PlayerX,
			&eGame.PlayerO,
			&eGame.CurrentTurn,
			&eGame.WinnerId,
			&eGame.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		eGame.Field = fieldJSON
		eGames = append(eGames, eGame)
	}
	var games []*model.Game
	for _, elem := range eGames {
		games = append(games, entity_mapper.GameToDomainFromEntity(&elem))
	}
	return games, nil
}

func (g gameRepository) GetFinishedByUserId(ctx context.Context, userId uuid.UUID) ([]*model.Game, error) {
	query := "SELECT id, field, status, player_x, player_o, winner_id, created_at FROM games WHERE status = 2 AND (player_x = $1 OR player_o = $1)"

	rows, err := g.db.DB.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("wrong query")
	}
	var eGames []entity.Game
	for rows.Next() {
		var eGame entity.Game
		var fieldJSON []byte

		err := rows.Scan(
			&eGame.ID,
			&fieldJSON,
			&eGame.Status,
			&eGame.PlayerX,
			&eGame.PlayerO,
			&eGame.WinnerId,
			&eGame.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		eGame.Field = fieldJSON
		eGames = append(eGames, eGame)
	}
	var games []*model.Game
	for _, elem := range eGames {
		games = append(games, entity_mapper.GameToDomainFromEntity(&elem))
	}
	return games, nil
}

func (g gameRepository) GetLeaderboard(ctx context.Context, count int) ([]*model.Leaderboard, error) {
	query := `
	WITH player_games AS (
		SELECT id AS game_id, player_x AS user_id, winner_id
		FROM games WHERE player_x IS NOT NULL AND games.status = 2
		UNION ALL
		SELECT id AS game_id, player_o AS user_id, winner_id
		FROM games WHERE player_o IS NOT NULL AND games.status = 2
	),
	user_stats AS (
		SELECT
			pg.user_id,
			COUNT(pg.game_id) AS total_games,
			SUM(CASE WHEN pg.user_id = pg.winner_id THEN 1 ELSE 0 END) AS wins
		FROM player_games pg
		GROUP BY pg.user_id
	)
	SELECT
		u.id AS user_id, u.login as login,
		ROUND((us.wins::numeric / NULLIF(us.total_games, 0)), 2) AS winrate
	FROM users u
	LEFT JOIN user_stats us ON u.id = us.user_id
	WHERE us.total_games > 0
	ORDER BY winrate DESC
	LIMIT $1;`

	rows, err := g.db.DB.Query(ctx, query, count)
	if err != nil {
		return nil, err
	}

	var eBoards []entity.Leaderboard
	for rows.Next() {
		var eBoard entity.Leaderboard
		err := rows.Scan(
			&eBoard.UserId,
			&eBoard.Login,
			&eBoard.Winrate)
		if err != nil {
			return nil, err
		}

		eBoards = append(eBoards, eBoard)
	}

	var mBoards []*model.Leaderboard

	for _, elem := range eBoards {
		mBoards = append(mBoards, entity_mapper.LeaderBoardToDomainFromEntity(&elem))
	}

	return mBoards, nil

}

func (g gameRepository) Update(ctx context.Context, game *entity.Game) error {
	query := "UPDATE games SET field = $2, status = $3, player_x = $4, player_o = $5, current_turn = $6, winner_id = $7 WHERE id = $1"
	_, err := g.db.DB.Exec(ctx, query, game.ID, game.Field, game.Status, game.PlayerX, game.PlayerO, game.CurrentTurn, game.WinnerId)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}

func (g gameRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM games where id = $1"
	_, err := g.db.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}
	return nil
}
