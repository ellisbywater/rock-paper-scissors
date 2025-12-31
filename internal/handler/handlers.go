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
	Name string `json:"name"`
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
