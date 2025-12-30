package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	ID int

	TotalRounds  int
	PlayerOne    Player
	PlayerTwo    Player
	CurrentRound int

	Rounds []Round
	Winner Player
	Score  Score
}

type NewGameRequest struct {
	TotalRounds int    `json:"total_rounds"`
	PlayerOne   Player `json:"player_one"`
	PlayerTwo   Player `json:"player_two"`
}

type NewPlayerRequest struct {
	Name string `json:"name"`
}

type RoundPlayerInput struct {
	PlayerID int  `json:"id"`
	Hand     Hand `json:"hand"`
}

type Round struct {
	Count     int
	PlayerOne RoundPlayerInput
	PlayerTwo RoundPlayerInput
	Winner    int
}

type RoundResult struct {
	RoundCount int
	Winner     int
}

type Score struct {
	PlayerOneScore int
	PlayerTwoScore int
}

type GameStore struct {
	games   []Game
	players []Player
}

type GameStoreService struct {
	Store GameStore
}

func (gs *GameStoreService) NewGame(req NewGameRequest) {
	id := len(gs.Store.games) + 1
	gs.Store.games = append(gs.Store.games, Game{ID: id, PlayerOne: req.PlayerOne, PlayerTwo: req.PlayerTwo, TotalRounds: req.TotalRounds, CurrentRound: 1})
}

func initGameStore() GameStoreService {
	return GameStoreService{
		Store: GameStore{
			games:   []Game{},
			players: []Player{},
		},
	}
}

func (gs *GameStoreService) NewRound(gameId int, playerOneID int, playerTwoID int) {
	game := gs.Store.games[gameId-1]
	if game.TotalRounds > len(game.Rounds) {
		newRound := Round{}
		game.Rounds = append(game.Rounds, newRound)
	}

}

func (gs *GameStoreService) NewPlayer(req NewPlayerRequest) {
	id := len(gs.Store.players) + 1
	player := Player{
		ID:   id,
		Name: req.Name,
	}
	gs.Store.players = append(gs.Store.players, player)

}

type Handlers struct {
	Service GameStoreService
}

func ResolveHands(hand_one Hand, hand_two Hand) int {
	switch hand_one {
	case Rock:
		if hand_two == Paper {
			return 2
		}
		if hand_two == Scissors {
			return 1
		}
		if hand_two == Rock {
			return 0
		}
	case Scissors:
		if hand_two == Rock {
			return 2
		}
		if hand_two == Paper {
			return 1
		}
		if hand_two == Scissors {
			return 0
		}
	case Paper:
		if hand_two == Scissors {
			return 2
		}
		if hand_two == Rock {
			return 1
		}
		if hand_two == Paper {
			return 0
		}
	}
	return 0
}

func (g *Game) RunRound(player_one_hand RoundPlayerInput, player_two_hand RoundPlayerInput) RoundResult {
	result := ResolveHands(player_one_hand.Hand, player_two_hand.Hand)
	var winner int = 0

	if result == 1 {
		winner = player_one_hand.PlayerID
	}
	if result == 2 {
		winner = player_two_hand.PlayerID
	}
	round := g.CurrentRound
	g.CurrentRound++

	return RoundResult{RoundCount: round, Winner: winner}
}

func (h *Handlers) NewGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
	var new_game NewGameRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&new_game); err != nil {
		http.Error(w, "incorrect game request", http.StatusBadRequest)
	}

	h.Service.NewGame(new_game)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) PlayRoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}

	// playerId, err := strconv.Atoi(r.PathValue("playerId"))
	// if err != nil {
	// 	http.Error(w, "Incorrect player id", http.StatusBadRequest)
	// }
	gameId, err := strconv.Atoi(r.PathValue("gameId"))
	if err != nil {
		http.Error(w, "Incorrect game id", http.StatusBadRequest)
	}
	roundNum, err := strconv.Atoi(r.PathValue("round"))
	if err != nil {
		http.Error(w, "Incorrect round id", http.StatusBadRequest)
	}

	var playerInput RoundPlayerInput

	game := h.Service.Store.games[gameId-1]
	round := game.Rounds[roundNum]

	if round.PlayerOne.PlayerID > 0 {
		if playerInput.PlayerID == round.PlayerOne.PlayerID {
			round.PlayerOne.Hand = playerInput.Hand
		}
	}
	if playerInput.PlayerID == round.PlayerTwo.PlayerID {
		round.PlayerTwo.Hand = playerInput.Hand
	}
}

func (h *Handlers) NewPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
	var new_player NewPlayerRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&new_player); err != nil {
		http.Error(w, "incorrect player request", http.StatusBadRequest)
	}

	h.Service.NewPlayer(new_player)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": fmt.Sprintf("User %s created", new_player.Name)})
}

type App struct {
	Port    string
	Service GameStoreService
}

func main() {
	app := App{
		Port:    ":8080",
		Service: initGameStore(),
	}
	handlers := Handlers{Service: app.Service}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /game/new", handlers.NewGameHandler)

	mux.HandleFunc("/player/new", handlers.NewPlayerHandler)

	log.Println("Rock Paper Scissors running on Port ", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, mux))
}
