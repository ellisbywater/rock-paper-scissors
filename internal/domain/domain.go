package domain

import (
	"context"
	"time"
)

/**

	Player
	- Has id
	- Belongs to many games
	- Belongs to many rounds
	- Belongs to many scores

	Game
	- Has Id
	- Has many rounds
	- Finished bool len(Rounds) + 1 = TotalRounds
	- Has many scores

	Rounds
	- Has Id
	- Has many Players
	- Has many Hands
**/

type Hand string

// const (
// 	Rock Hand = iota
// 	Paper
// 	Scissors
// )

// func (h Hand) String() string {
// 	return [...]string{"rock", "paper", "scissors"}[h]
// }

type PlayerCreateRequest struct {
	UserName string `json:"username"`
}

type PlayerResponse struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type Game struct {
	ID          int
	CreatedAt   time.Time
	TotalRounds int

	PlayerOne      PlayerResponse
	PlayerTwo      PlayerResponse
	PlayerOneScore int
	PlayerTwoScore int

	Rounds []Round
	Winner int
}

type GameResponse struct {
	ID             int       `json:"id"`
	TotalRounds    int       `json:"total_rounds"`
	CurrentRound   int       `json:"current_round"`
	PlayerOneId    int       `json:"player_one_id"`
	PlayerTwoId    int       `json:"player_two_id"`
	PlayerOneScore int       `json:"player_one_score"`
	PlayerTwoScore int       `json:"player_two_score"`
	Winner         int       `json:"winner"`
	Finished       bool      `json:"finished"`
	Rounds         []Round   `json:"rounds"`
	CreatedAt      time.Time `json:"created_at"`
}

type GameCreateResponse struct {
	ID           int       `json:"id"`
	TotalRounds  int       `json:"total_rounds"`
	CurrentRound int       `json:"current_round"`
	PlayerOneId  int       `json:"player_one_id"`
	PlayerTwoId  int       `json:"player_two_id"`
	CreatedAt    time.Time `json:"created_at"`
}
type GameCreateRequest struct {
	TotalRounds int `json:"total_rounds"`
	PlayerOneID int `json:"player_one_id"`
	PlayerTwoID int `json:"player_two_id"`
}

type Round struct {
	ID            int
	GameId        int
	Count         int
	PlayerOneID   int
	PlayerTwoID   int
	PlayerOneHand string
	PlayerTwoHand string
	Winner        int
	Finished      bool
}

type RoundCreateRequest struct {
	GameId int `json:"game_id"`
}

type RoundCreateResponse struct {
	Id          int `json:"id"`
	GameId      int `json:"game_id"`
	Count       int `json:"count"`
	PlayerOneID int `json:"player_one_id"`
	PlayerTwoID int `json:"player_two_id"`
}

type PlayerScore struct {
	PlayerID int `json:"player_id"`
	Score    int `json:"score"`
}

type Score struct {
	PlayerOne PlayerScore `json:"player_one_score"`
	PlayerTwo PlayerScore `json:"player_two_score"`
}

type RoundPlayerInput struct {
	RoundId  int    `json:"round_id"`
	GameID   int    `json:"game_id"`
	PlayerID int    `json:"player_id"`
	Hand     string `json:"hand"`
}

type PlayerRepository interface {
	Create(ctx context.Context, player PlayerCreateRequest, res *PlayerResponse) error
	Get(ctx context.Context, id int, res *PlayerResponse) error
	GetGames(ctx context.Context, id int, res *[]GameResponse) error
}

type GameRepository interface {
	Create(ctx context.Context, game GameCreateRequest, res *GameCreateResponse) error
	Get(ctx context.Context, id int, res *GameResponse) error
}

type RoundRepository interface {
	Create(ctx context.Context, round_create_request RoundCreateRequest, res *RoundCreateResponse) error
	UpdateHand(ctx context.Context, player_input RoundPlayerInput, res *Round) error
	Get(ctx context.Context, id int, res *Round) error
}
