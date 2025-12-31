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

	Rounds []RoundContext
	Winner int
}

type GameResponse struct {
	ID             int            `json:"id"`
	TotalRounds    int            `json:"total_rounds"`
	CurrentRound   int            `json:"current_round"`
	PlayerOneId    int            `json:"player_one_id"`
	PlayerTwoId    int            `json:"player_two_id"`
	PlayerOneScore int            `json:"player_one_score"`
	PlayerTwoScore int            `json:"player_two_score"`
	Winner         int            `json:"winner"`
	Finished       bool           `json:"finished"`
	Rounds         []RoundContext `json:"rounds"`
	CreatedAt      time.Time      `json:"created_at"`
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

type RoundContext struct {
	ID            int    `json:"id"`
	GameID        int    `json:"game_id"`
	Count         int    `json:"count"`
	PlayerOneID   int    `json:"player_one_id"`
	PlayerTwoID   int    `json:"player_two_id"`
	CurrentPlayer int    `json:"current_player"`
	PlayerOneHand string `json:"player_one_hand"`
	PlayerTwoHand string `json:"player_two_hand"`
	Winner        int    `json:"winner"`
	Finished      bool   `json:"finished"`
}

type PlayerHandContext struct {
	ID   int
	Hand string
}

func (rc *RoundContext) PlayerOneHandContext() PlayerHandContext {
	ctx := PlayerHandContext{
		ID:   rc.PlayerOneID,
		Hand: rc.PlayerOneHand,
	}
	return ctx
}

func (rc *RoundContext) PlayerTwoHandContext() PlayerHandContext {
	ctx := PlayerHandContext{
		ID:   rc.PlayerTwoID,
		Hand: rc.PlayerTwoHand,
	}
	return ctx
}

func (rc *RoundContext) HasPlayerOnePlayed() bool {
	if rc.PlayerOneHand == "none" {
		return false
	}
	return true
}

func (rc *RoundContext) HasPlayerTwoPlayed() bool {
	if rc.PlayerTwoHand == "none" {
		return false
	}
	return true
}

func (rc *RoundContext) CurrentPlayerHand() string {
	switch rc.CurrentPlayer {
	case rc.PlayerOneID:
		return rc.PlayerOneHand
	case rc.PlayerTwoID:
		return rc.PlayerTwoHand
	default:
		return ""
	}
}

type WinnerContext struct {
	RoundID  int
	Hand     string
	PlayerID int
}

func (rc *RoundContext) CalculateWinner() WinnerContext {
	var winnerCtx WinnerContext
	winnerCtx.RoundID = rc.ID
	if rc.PlayerOneHand != "none" && rc.PlayerTwoHand != "none" {
		switch rc.PlayerOneHand {
		case "rock":
			if rc.PlayerTwoHand == "paper" {
				winnerCtx.PlayerID = rc.PlayerTwoID
				winnerCtx.Hand = rc.PlayerTwoHand
			}
			if rc.PlayerTwoHand == "scissors" {
				winnerCtx.PlayerID = rc.PlayerOneID
				winnerCtx.Hand = rc.PlayerOneHand
			}
			if rc.PlayerTwoHand == "rock" {
				winnerCtx.PlayerID = 0
			}
		case "scissors":
			if rc.PlayerTwoHand == "rock" {
				winnerCtx.PlayerID = rc.PlayerTwoID
				winnerCtx.Hand = rc.PlayerTwoHand
			}
			if rc.PlayerTwoHand == "paper" {
				winnerCtx.PlayerID = rc.PlayerOneID
				winnerCtx.Hand = rc.PlayerOneHand
			}
			if rc.PlayerTwoHand == "scissors" {
				winnerCtx.PlayerID = 0
			}
		case "paper":
			if rc.PlayerTwoHand == "scissors" {
				winnerCtx.PlayerID = rc.PlayerTwoID
				winnerCtx.Hand = rc.PlayerTwoHand
			}
			if rc.PlayerTwoHand == "rock" {
				winnerCtx.PlayerID = rc.PlayerOneID
				winnerCtx.Hand = rc.PlayerOneHand
			}
			if rc.PlayerTwoHand == "paper" {
				winnerCtx.PlayerID = 0
			}
		}
	}
	return winnerCtx
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
	Create(ctx context.Context, round_create_request RoundContext, res *RoundContext) error
	UpdateHand(ctx context.Context, player_input RoundContext, res *RoundContext) error
	Get(ctx context.Context, id int, res *RoundContext) error
}
