package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/service"
)

type RoundResultRequest struct {
	RoundCount int `json:"id"`
	Winner     int `json:"winner"`
}

type NewGameRequest struct {
	TotalRounds int `json:"total_rounds"`
	PlayerOne   int `json:"player_one"`
	PlayerTwo   int `json:"player_two"`
}

type NewPlayerRequest struct {
	UserName string `json:"username"`
}

type RoundPlayerInput struct {
	PlayerID int `json:"id"`
	Hand     int `json:"hand"`
}

type GameHandlers struct {
	service service.GameService
}

func NewGameHandler(service service.GameService) *GameHandlers {
	return &GameHandlers{service: service}
}

func (gh *GameHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var new_game_req NewGameRequest
	if err := json.NewDecoder(r.Body).Decode(&new_game_req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	defer r.Body.Close()
	if new_game_req.TotalRounds < 1 {
		new_game_req.TotalRounds = 1
	}
	game, err := gh.service.NewGame(r.Context(), new_game_req.TotalRounds, new_game_req.PlayerOne, new_game_req.PlayerTwo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

func (gh *GameHandlers) GetGame(w http.ResponseWriter, r *http.Request) {
	game_id, err := strconv.Atoi(r.PathValue("gameId"))
	if err != nil {
		http.Error(w, "Invalid Game ID", http.StatusBadRequest)
	}
	game, err := gh.service.GetGame(r.Context(), game_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(game)
}

type PlayerHandlers struct {
	service service.PlayerService
}

func NewPlayerHandler(service service.PlayerService) *PlayerHandlers {
	return &PlayerHandlers{service: service}
}

func (ph *PlayerHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var new_player_req NewPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&new_player_req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	defer r.Body.Close()
	if new_player_req.UserName == "" {
		http.Error(w, "Username cannot be blank", http.StatusBadRequest)
	}
	player, err := ph.service.CreatePlayer(r.Context(), new_player_req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

func (ph *PlayerHandlers) Get(w http.ResponseWriter, r *http.Request) {
	player_id, err := strconv.Atoi(r.PathValue("playerId"))
	if err != nil {
		http.Error(w, "Invalid player id", http.StatusBadRequest)
	}
	player, err := ph.service.GetPlayer(r.Context(), player_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(player)
}

func (ph *PlayerHandlers) GetGames(w http.ResponseWriter, r *http.Request) {
	player_id, err := strconv.Atoi(r.PathValue("playerId"))
	if err != nil {
		http.Error(w, "Invalid player id", http.StatusBadRequest)
	}
	games, err := ph.service.GetPlayerGames(r.Context(), player_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(games)
}
