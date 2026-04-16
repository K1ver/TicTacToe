package service

import (
	"app/internal/datasource/entity"
	"app/internal/datasource/entity_mapper"
	"app/internal/datasource/repository"
	"app/internal/domain/algorithm"
	"app/internal/domain/model"
	"app/internal/web/dto"
	"app/internal/web/dto_mapper"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

const botId = "00000000-0000-0000-0000-000000000001"

const (
	PlayerX = 1
	PlayerO = 2
)

type GameService interface {
	// CreateGame создаёт новую игру принимая размер поля, сторону и выбор с кем играть
	CreateGame(ctx context.Context, rows, cols, side int, isCPU bool, playerId uuid.UUID) *dto.Game
	// JoinGame присоединение к игре
	JoinGame(ctx context.Context, gameId uuid.UUID, playerId uuid.UUID) (*dto.Game, error)
	// MakeMove сделать ход
	MakeMove(ctx context.Context, gameId uuid.UUID, playerId uuid.UUID, newGF model.GameField) (*dto.Game, error)
	// ValidateGame проверка, что поле не было некорректно изменено
	ValidateGame(ctx context.Context, gameId uuid.UUID, newGF model.GameField) (bool, error)
	// EndGame проверка на конец игры и победителя
	EndGame(ctx context.Context, gameId uuid.UUID) (*dto.Game, error)
	// GetGamesByStatus получение всех игр по статусу
	GetGamesByStatus(ctx context.Context, status model.GameStatus) ([]*dto.Game, error)
	// GetFinishedGamesByUserId получения всех завершенных игр по UUID пользователя.
	GetFinishedGamesByUserId(ctx context.Context, userId uuid.UUID) ([]*dto.Game, error)
	//GetLeaderboard получение N записей пользователей и соотношение побед и поражений
	GetLeaderboard(ctx context.Context, count int) ([]*dto.Leaderboard, error)
}

type gameService struct {
	repo      repository.GameRepository
	algorithm algorithm.MinimaxAlgorithm
}

func NewGameService(repo repository.GameRepository, algorithm algorithm.MinimaxAlgorithm) GameService {
	return &gameService{
		repo:      repo,
		algorithm: algorithm,
	}
}

func (g *gameService) CreateGame(ctx context.Context, rows, cols, side int, isCPU bool, player uuid.UUID) *dto.Game {
	newGame := model.NewGame(rows, cols)
	if side == PlayerO {
		newGame.PlayerO = player
	} else {
		newGame.PlayerX = player
	}

	if isCPU {
		newGame.Status = model.StatusInProgress
		if newGame.PlayerX == uuid.Nil {
			newGame.PlayerX = uuid.MustParse(botId)
			newGame.CurrentTurn = newGame.PlayerO
			g.makeNextMoveByMinimax(&newGame, 4)
		} else {
			newGame.PlayerO = uuid.MustParse(botId)
			newGame.CurrentTurn = newGame.PlayerX
		}

	}

	err := g.repo.Create(ctx, entity_mapper.GameToEntity(&newGame))
	if err != nil {
		log.Print(err)
	}
	return dto_mapper.GameToDto(&newGame)
}

func (g *gameService) JoinGame(ctx context.Context, gameId uuid.UUID, playerId uuid.UUID) (*dto.Game, error) {
	game, err := g.repo.GetByID(ctx, gameId)
	if err != nil {
		return nil, err
	}
	if game.Status != model.StatusIsWaiting {
		return nil, fmt.Errorf("the game has already started")
	}

	if game.PlayerX == playerId || game.PlayerO == playerId {
		return nil, fmt.Errorf("you're the creator you cannot join as a second player")
	}

	if game.PlayerX == uuid.Nil {
		game.PlayerX = playerId
	} else {
		game.PlayerO = playerId
	}
	game.Status = model.StatusInProgress
	game.CurrentTurn = game.PlayerX
	err = g.repo.Update(ctx, entity_mapper.GameToEntity(game))
	if err != nil {
		return nil, err
	}
	return dto_mapper.GameToDto(game), nil
}

func (g *gameService) ValidateGame(ctx context.Context, uuid uuid.UUID, newGF model.GameField) (bool, error) {
	prevGameMove, err := g.repo.GetByID(ctx, uuid)
	if err != nil {
		return false, err
	}

	if model.StatusIsWaiting == prevGameMove.Status {
		return false, fmt.Errorf("the game has't yet started")
	}

	if model.StatusFinished == prevGameMove.Status {
		return false, fmt.Errorf("the game has already finished")
	}

	rows, cols := newGF.GetSize()
	prevRows, prevCols := prevGameMove.Field.GetSize()

	if (prevRows != rows) || (prevCols != cols) {
		return false, nil
	}

	changes := 0

	for i := range rows {
		for j := range cols {
			if (prevGameMove.Field[i][j] == 1 && newGF[i][j] != 1) || (prevGameMove.Field[i][j] == 2 && newGF[i][j] != 2) {
				return false, nil
			}
			if prevGameMove.Field[i][j] != newGF[i][j] {
				changes++
			}
		}
	}

	return changes == 1, nil
}

func (g *gameService) EndGame(ctx context.Context, uuid uuid.UUID) (*dto.Game, error) {
	game, err := g.repo.GetByID(ctx, uuid)
	if err != nil {
		return &dto.Game{}, err
	}

	if g.checkWinner(game.Field, PlayerX) {
		game.Status = model.StatusFinished
		game.WinnerId = game.PlayerX
	} else if g.checkWinner(game.Field, PlayerO) {
		game.Status = model.StatusFinished
		game.WinnerId = game.PlayerO
	} else if g.isBoardFull(game.Field) {
		game.Status = model.StatusFinished
	}

	_ = g.repo.Update(ctx, entity_mapper.GameToEntity(game))
	return dto_mapper.GameToDto(game), nil
}

func (g *gameService) MakeMove(ctx context.Context, gameId uuid.UUID, playerId uuid.UUID, newGF model.GameField) (*dto.Game, error) {
	game, err := g.repo.GetByID(ctx, gameId)
	if err != nil {
		return &dto.Game{}, err
	}
	if game.CurrentTurn != playerId {
		return nil, fmt.Errorf("other player turn")
	}
	game.Field = newGF
	if g.isBoardFull(game.Field) {
		game.CurrentTurn = uuid.Nil
		_ = g.repo.Update(ctx, entity_mapper.GameToEntity(game))
		return dto_mapper.GameToDto(game), nil
	}

	if game.Status == model.StatusInProgress && (game.PlayerO == uuid.MustParse(botId) || game.PlayerX == uuid.MustParse(botId)) {
		g.makeNextMoveByMinimax(game, 4)
	} else {
		if game.CurrentTurn == game.PlayerX {
			game.CurrentTurn = game.PlayerO
		} else {
			game.CurrentTurn = game.PlayerX
		}
	}
	_ = g.repo.Update(ctx, entity_mapper.GameToEntity(game))
	return dto_mapper.GameToDto(game), nil
}

func (g *gameService) GetGamesByStatus(ctx context.Context, status model.GameStatus) ([]*dto.Game, error) {
	mGames, err := g.repo.GetAllByStatus(ctx, entity.GameStatus(status))
	if err != nil {
		return nil, err
	}

	var dGames []*dto.Game

	for _, elem := range mGames {
		dGames = append(dGames, dto_mapper.GameToDto(elem))
	}
	return dGames, nil
}

func (g *gameService) GetFinishedGamesByUserId(ctx context.Context, userId uuid.UUID) ([]*dto.Game, error) {
	mGames, err := g.repo.GetFinishedByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	var dGames []*dto.Game

	for _, elem := range mGames {
		dGames = append(dGames, dto_mapper.GameToDto(elem))
	}
	return dGames, nil
}

func (g *gameService) GetLeaderboard(ctx context.Context, count int) ([]*dto.Leaderboard, error) {
	mBoards, err := g.repo.GetLeaderboard(ctx, count)
	if err != nil {
		return nil, err
	}

	var dBoards []*dto.Leaderboard

	for _, elem := range mBoards {
		dBoards = append(dBoards, dto_mapper.LeaderBoardToDto(elem))
	}

	return dBoards, nil
}

func (g *gameService) makeNextMoveByMinimax(game *model.Game, difficult int) {
	//TODO сделать сложность
	/*
		Сложно = 6
		Средне = 4
		Легко = 2
	*/
	var row, col int
	if game.PlayerX == uuid.MustParse(botId) {
		row, col, _ = g.algorithm.FindBestMove(game.Field, PlayerX, difficult)
		game.Field[row][col] = PlayerX
	} else {
		row, col, _ = g.algorithm.FindBestMove(game.Field, PlayerO, difficult)
		game.Field[row][col] = PlayerO
	}
}

func (g *gameService) checkWinner(field model.GameField, player int) bool {
	rows, cols := field.GetSize()

	// Проверка строк и столбов
	for i := 0; i < rows; i++ {
		for j := 0; j <= cols-3; j++ {
			if field[i][j] == player && field[i][j+1] == player && field[i][j+2] == player {
				return true
			}
		}
	}

	for j := 0; j < cols; j++ {
		for i := 0; i <= rows-3; i++ {
			if field[i][j] == player && field[i+1][j] == player && field[i+2][j] == player {
				return true
			}
		}
	}

	// Проверка диагоналей (сверху-слева направо-вниз)
	for i := 0; i <= rows-3; i++ {
		for j := 0; j <= cols-3; j++ {
			if field[i][j] == player && field[i+1][j+1] == player && field[i+2][j+2] == player {
				return true
			}
		}
	}

	// Проверка диагоналей (сверху-справа налево-вниз)
	for i := 0; i <= rows-3; i++ {
		for j := 2; j < cols; j++ {
			if field[i][j] == player && field[i+1][j-1] == player && field[i+2][j-2] == player {
				return true
			}
		}
	}

	return false
}

func (g *gameService) isBoardFull(field model.GameField) bool {
	rows, cols := field.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if field[i][j] == 0 {
				return false
			}
		}
	}
	return true
}
