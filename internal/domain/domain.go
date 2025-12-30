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

type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

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
	PlayerOne   PlayerResponse
	PlayerTwo   PlayerResponse

	Rounds []Round
	Winner int
	Score  Score
}

type GameResponse struct {
	ID          int            `json:"id"`
	TotalRounds int            `json:"total_rounds"`
	PlayerOne   PlayerResponse `json:"player_one"`
	PlayerTwo   PlayerResponse `json:"player_two"`
	Winner      int            `json:"winner"`
	Score       Score          `json:"score"`
	Finished    bool           `json:"finished"`
}

type GameCreateRequest struct {
	TotalRounds int `json:"total_rounds"`
	PlayerOneID int `json:"player_one_id"`
	PlayerTwoID int `json:"player_two_id"`
}

type Round struct {
	ID        int
	Count     int
	PlayerOne RoundPlayerInput
	PlayerTwo RoundPlayerInput
	Winner    int
}

type RoundCreateRequest struct {
	Count       int `json:"count"`
	PlayerOneId int `json:"player_one_id"`
	PlayerTwoID int `json:"player_two_id"`
}

type RoundCreateResponse struct {
	Id          int `json:"id"`
	PlayerOneId int `json:"player_one_id"`
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
	PlayerID int  `json:"player_id"`
	Hand     Hand `json:"hand"`
}

type PlayerRepository interface {
	Create(ctx context.Context, player PlayerCreateRequest) (error, PlayerResponse)
	Get(ctx context.Context, id int) (error, PlayerResponse)
	GetGames(ctx context.Context, id int) (error, []GameResponse)
}

type GameRepository interface {
	Create(ctx context.Context, game GameCreateRequest) (error, GameResponse)
	Get(ctx context.Context, id int) (error, GameResponse)
}

type RoundRepository interface {
	Create(ctx context.Context, round_create_request RoundCreateRequest) (error, RoundCreateResponse)
	PlayHand(ctx context.Context, player_input RoundPlayerInput)
	Get(ctx context.Context, id int) (error, RoundCreateResponse)
}
