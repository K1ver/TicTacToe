package handler

import (
	"app/internal/domain/model"
	"app/internal/domain/service"
	"app/internal/web/dto"
	"app/internal/web/dto_mapper"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GameHandler interface {
	// MakeMoveHandle обрабатывает ход игрока
	MakeMoveHandle(w http.ResponseWriter, r *http.Request)
	// CreateGameHandle создает новую игру
	CreateGameHandle(w http.ResponseWriter, r *http.Request)
	// JoinGameHandle присоединение к игре
	JoinGameHandle(w http.ResponseWriter, r *http.Request)
	// GetGameHandle получает игру по id
	GetGameHandle(w http.ResponseWriter, r *http.Request)
	// GetAllUnstartedGamesHandle браузер игр
	GetAllUnstartedGamesHandle(w http.ResponseWriter, r *http.Request)
	// GetAllFinishedGamesHandle получение личных завершенных игр
	GetAllFinishedGamesHandle(w http.ResponseWriter, r *http.Request)
	//GetLeaderboardHandle получение N записей пользователей и соотношение побед и поражений
	GetLeaderboardHandle(w http.ResponseWriter, r *http.Request)
}

type gameHandle struct {
	gameService service.GameService
	userService service.UserService
}

func NewGameHandler(gameService service.GameService, userService service.UserService) GameHandler {
	return &gameHandle{
		gameService: gameService,
		userService: userService,
	}
}

func (g *gameHandle) MakeMoveHandle(w http.ResponseWriter, r *http.Request) {
	// CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	gameIdStr := chi.URLParam(r, "gameId")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	var req dto.Game
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	domainReq := dto_mapper.GameToDomainFromDto(&req)
	val, err := g.gameService.ValidateGame(ctx, gameId, domainReq.Field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	if !val {
		http.Error(w, "Invalid Move", http.StatusBadRequest)
		log.Println(err)
		return
	}
	user := g.getUserByHeader(ctx, r)
	_, err = g.gameService.MakeMove(ctx, gameId, user.Id, domainReq.Field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	response, err := g.gameService.EndGame(ctx, gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (g *gameHandle) CreateGameHandle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	type settingGame struct {
		Side        int  `json:"side"`
		GameWithCPU bool `json:"gameWithCPU"`
	}

	var s settingGame

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Wrong body", http.StatusBadRequest)
		return
	}

	if s.Side == 0 {
		s.Side = 1
	}

	user := g.getUserByHeader(ctx, r)

	dtoGame := g.gameService.CreateGame(ctx, 3, 3, s.Side, s.GameWithCPU, user.Id)
	response := struct {
		Id string `json:"id"`
	}{
		Id: dtoGame.Id.String(),
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (g *gameHandle) JoinGameHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	gameIdStr := chi.URLParam(r, "gameId")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := g.getUserByHeader(ctx, r)

	playerId := user.Id
	game, err := g.gameService.JoinGame(ctx, gameId, playerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (g *gameHandle) GetGameHandle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	gameIdStr := chi.URLParam(r, "gameId")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dtoGame, err := g.gameService.EndGame(ctx, gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(dtoGame); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (g *gameHandle) GetAllUnstartedGamesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	games, err := g.gameService.GetGamesByStatus(ctx, model.StatusIsWaiting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type response struct {
		Id   uuid.UUID `json:"id"`
		Side int       `json:"side"`
	}

	var res []response

	for _, elem := range games {
		var r response
		r.Id = elem.Id
		if elem.PlayerX == uuid.Nil {
			r.Side = 1
		} else {
			r.Side = 2
		}

		res = append(res, r)
	}

	if len(res) == 0 {
		_ = json.NewEncoder(w).Encode([]response{})
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (g *gameHandle) GetAllFinishedGamesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	user := g.getUserByHeader(ctx, r)
	log.Print(user.Id)
	games, err := g.gameService.GetFinishedGamesByUserId(ctx, user.Id)
	log.Print(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(games); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func (g *gameHandle) getUserByHeader(ctx context.Context, r *http.Request) *dto.User {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	player, _ := g.userService.GetUserByToken(ctx, tokenString)
	return player
}

func (g *gameHandle) GetLeaderboardHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	countStr := chi.URLParam(r, "count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board, err := g.gameService.GetLeaderboard(ctx, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(board) == 0 {
		_ = json.NewEncoder(w).Encode([]dto.Leaderboard{})
		return
	}

	if err = json.NewEncoder(w).Encode(board); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
