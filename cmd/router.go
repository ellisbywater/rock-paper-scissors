package main

import (
	"database/sql"
	"net/http"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/handler"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/repository"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/service"
)

func buildGameHandlerDeps(db *sql.DB) handler.GameHandlers {
	var gameRepo domain.GameRepository = repository.NewGameRepository(db)
	var gameService service.GameService = *service.NewGameService(gameRepo)
	return *handler.NewGameHandler(gameService)
}

func buildPlayerHandlerDeps(db *sql.DB) handler.PlayerHandlers {
	var playerRepo domain.PlayerRepository = repository.NewPlayerRepository(db)
	var playerService service.PlayerService = *service.NewPlayerService(playerRepo)
	return *handler.NewPlayerHandler(playerService)
}

func buildRoundHandlerDeps(db *sql.DB) handler.RoundHandlers {
	var roundRepo domain.RoundRepository = repository.NewRoundRepository(db)
	var roundService service.RoundService = *service.NewRoundService(roundRepo)
	return *handler.NewRoundHandlers(roundService)
}

func NewRouterWithDeps(db *sql.DB) *http.ServeMux {
	r := http.NewServeMux()

	gameHandler := buildGameHandlerDeps(db)
	playerHandler := buildPlayerHandlerDeps(db)
	roundHandler := buildRoundHandlerDeps(db)

	r.HandleFunc("POST /player/create", playerHandler.Create)
	r.HandleFunc("GET /player/{playerId}", playerHandler.Get)
	r.HandleFunc("GET /player/{playerId}/games", playerHandler.GetGames)

	r.HandleFunc("POST /game/create", gameHandler.Create)
	r.HandleFunc("GET /game/{gameId}", gameHandler.GetGame)

	r.HandleFunc("POST /round/create", roundHandler.Create)

	return r
}
