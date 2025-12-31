package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/handler"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/repository"
	"github.com/ellisbywater/http-rock-paper-scissors/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

type App struct {
	Port   string
	Router *http.ServeMux
}

const port = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("error in connection with database %s", err)
	}

	r := http.NewServeMux()

	gameHandler := buildGameHandlerDeps(db)
	playerHandler := buildPlayerHandlerDeps(db)
	roundHandler := buildRoundHandlerDeps(db)

	r.HandleFunc("POST /player/create", playerHandler.Create)
	r.HandleFunc("GET /player/{playerId}", playerHandler.Get)
	r.HandleFunc("GET /player/{playerId}/games", playerHandler.GetGames)

	r.HandleFunc("POST /game/create", gameHandler.Create)
	r.HandleFunc("GET /game/{gameId}", gameHandler.GetGame)

	r.HandleFunc("POST /game/{gameId}/round/create", roundHandler.Create)
	r.HandleFunc("POST /game/{gameId}/round/{roundId}/playHand", roundHandler.PlayHand)

	log.Println("Connected to database")

	log.Println("Rock Paper Scissors running on Port ", port)
	log.Fatal(http.ListenAndServe(port, r))
}
