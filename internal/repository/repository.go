package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
)

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) domain.GameRepository {
	return &gameRepository{db}
}

func (gr *gameRepository) Create(ctx context.Context, game domain.GameCreateRequest, res *domain.GameCreateResponse) error {
	query := `
		INSERT INTO games (
			total_rounds,
			current_round,
			player_one_id,
			player_two_id
		) Values (
		 	$1,
			1,
			$2,
			$3
		 )
		RETURNING id, total_rounds, current_round, created_at, player_one_id, player_two_id;
	`
	err := gr.db.QueryRowContext(
		ctx,
		query,
		game.TotalRounds,
		game.PlayerOneID,
		game.PlayerTwoID,
	).Scan(&res.ID, &res.TotalRounds, &res.CurrentRound, &res.CreatedAt, &res.PlayerOneId, &res.PlayerTwoId)

	if err != nil {
		return err
	}
	return nil
}

func (gr *gameRepository) Get(ctx context.Context, id int, res *domain.GameResponse) error {
	// TODO: update query to join rounds
	query := `
		SELECT id, total_rounds, current_round, player_one_id, player_two_id, COALESCE(winner, 0), created_at
		FROM games
		WHERE id = $1;
	`
	err := gr.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.TotalRounds,
		&res.CurrentRound,
		&res.PlayerOneId,
		&res.PlayerTwoId,
		&res.Winner,
		&res.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) domain.PlayerRepository {
	return &playerRepository{db}
}

func (pr *playerRepository) Create(ctx context.Context, player domain.PlayerCreateRequest, res *domain.PlayerResponse) error {
	query := `
		INSERT INTO players (
			username
		) VALUES (
		 	$1
		) RETURNING id, username
	`
	err := pr.db.QueryRowContext(ctx, query, player.UserName).Scan(
		&res.ID,
		&res.UserName,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pr *playerRepository) Get(ctx context.Context, id int, res *domain.PlayerResponse) error {
	query := `
		SELECT * FROM players WHERE id=$1;
	`
	err := pr.db.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (pr *playerRepository) GetGames(ctx context.Context, id int, res *[]domain.GameResponse) error {
	query := `
		SELECT * FROM games AS g JOIN rounds AS r ON g.id = r.game WHERE g.player_one_id = $1 OR g.player_two_id = $1
	`
	rows, err := pr.db.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var game domain.GameResponse
		err := rows.Scan(
			&game.ID,
			&game.TotalRounds,
			&game.PlayerOneId,
			&game.PlayerTwoId,
			&game.PlayerOneScore,
			&game.PlayerTwoScore,
			&game.Rounds,
			&game.Finished,
			&game.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

type roundRepository struct {
	db *sql.DB
}

func NewRoundRepository(db *sql.DB) domain.RoundRepository {
	return &roundRepository{db}
}

func (rr *roundRepository) Get(ctx context.Context, id int, res *domain.RoundContext) error {
	query := `
		SELECT * FROM rounds WHERE id=$1;
	`
	err := rr.db.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.Count, &res.PlayerOneHand, &res.PlayerTwoHand)
	if err != nil {
		return err
	}
	return nil
}

func (rr *roundRepository) Create(ctx context.Context, round_create_request domain.RoundContext, res *domain.RoundContext) error {
	type gameContext struct {
		current_round int
		total_rounds  int
		finished      bool
		player_one_id int
		player_two_id int
	}
	var newGameContext gameContext
	check_count_query := `
		SELECT current_round, total_rounds, player_one_id, player_two_id, finished FROM games WHERE id=$1;
	`
	err := rr.db.QueryRowContext(ctx, check_count_query, round_create_request.GameID).Scan(&newGameContext.current_round, &newGameContext.total_rounds, &newGameContext.player_one_id, &newGameContext.player_two_id, &newGameContext.finished)
	if err != nil {
		return err
	}
	if newGameContext.finished {
		return errors.New("Game Already Finished!")
	}
	if newGameContext.current_round == newGameContext.total_rounds {
		return fmt.Errorf("Already on last round: %d", newGameContext.current_round)
	}
	query := `
		INSERT INTO rounds (
			game,
			count,
			player_one_id,
			player_two_id
		) VALUES (
		 	$1,
			$2,
			$3,
			$4
		) RETURNING id, game, count, player_one_id, player_two_id;
	`
	err = rr.db.QueryRowContext(
		ctx,
		query,
		round_create_request.GameID,
		newGameContext.current_round+1,
		newGameContext.player_one_id,
		newGameContext.player_two_id,
	).Scan(&res.ID, &res.GameID, &res.Count, &res.PlayerOneID, &res.PlayerTwoID)

	if err != nil {
		return err
	}
	return nil
}

// TODO: Has to be someway to simplify or condense the queries
// Checks For Winner
// Updates Score
// Updates game finished
func (rr *roundRepository) CheckForWinner(ctx context.Context, res *domain.RoundContext) error {
	type handsResults struct {
		player_one_id   int
		player_two_id   int
		player_one_hand domain.Hand
		player_two_hand domain.Hand
	}
	var results handsResults
	query_player_select := `
		SELECT player_one_id, player_two_id, COALESCE(player_one_hand, 'none'), COALESCE(player_two_hand, 'none') FROM rounds WHERE id = $1
	`
	// Retrieve Fields for comparison
	err := rr.db.QueryRowContext(ctx, query_player_select, res.ID).Scan(&results.player_one_id, &results.player_two_id, &results.player_one_hand, &results.player_two_hand)
	fmt.Println("CheckWinner HandsResults: ", results)
	if err != nil {
		return err
	}

	winner_query := `UPDATE rounds SET winner = $1, finished = True WHERE id = $2 RETURNING * `

	var winPlayerId int
	var player_score_name string //dynamic query string
	// Calculate score and create respective queries
	if results.player_one_hand != "none" && results.player_two_hand != "none" {
		winner_num := resolveHands(string(results.player_one_hand), string(results.player_two_hand))
		fmt.Println("Winner Number: ", winner_num)
		switch winner_num {
		case 1:
			winPlayerId = results.player_one_id
			player_score_name = "player_one_score"
		case 2:
			winPlayerId = results.player_two_id
			player_score_name = "player_two_score"
		default:
			winner_query = `UPDATE rounds SET winner=NULL, finished=True WHERE id=$1 RETURNING *;`
			player_score_name = ""
		}
	} else {
		// query for no winner (sad); write to res
		no_winner_query := `SELECT
							id,
							game, 
							count, 
							player_one_id, 
							player_two_id, 
							COALESCE(player_one_hand, 'none'), 
							COALESCE(player_two_hand, 'none'), 
							COALESCE(winner, 0), 
							finished
							FROM rounds WHERE id=$1`
		err := rr.db.QueryRowContext(ctx, no_winner_query, res.ID).Scan(&res.ID, &res.GameID, &res.Count, &res.PlayerOneID, &res.PlayerTwoID, &res.PlayerOneHand, &res.PlayerTwoHand, &res.Winner, &res.Finished)
		if err != nil {
			return err
		}
		return nil
	}

	// Update Round Winner
	row := rr.db.QueryRowContext(ctx, winner_query, winPlayerId, res.ID)
	fmt.Println("Update round winner row context: ", row)

	err = row.Scan(&res.ID, &res.GameID, &res.Count, &res.PlayerOneHand, &res.PlayerTwoHand, &res.Winner, &res.Finished)
	if err != nil {
		return err
	}

	type checkGameCount struct {
		total_rounds  int
		current_round int
	}
	var gameCount checkGameCount
	// Update current round
	err = rr.db.QueryRowContext(ctx, "UPDATE games SET current_round = current_round + 1 WHERE id=$1 RETURNING current_round, total_rounds;", res.GameID).Scan(&gameCount.current_round, &gameCount.total_rounds)
	if err != nil {
		return err
	}

	// Update Score
	var score_query string
	if len(player_score_name) != 0 {
		switch player_score_name {
		case "player_one_score":
			score_query = `UPDATE games SET player_one_score = player_one_score + 1 WHERE id=$1;`
		case "player_two_score":
			score_query = `UPDATE games SET player_two_score = player_two_score + 1 WHERE id=$1;`
		default:
			score_query = ""
		}
		if len(score_query) > 0 {
			_ = rr.db.QueryRowContext(ctx, score_query, res.GameID)
		}
	}

	if gameCount.current_round == gameCount.total_rounds {
		_ = rr.db.QueryRowContext(ctx, "UPDATE games SET finished=True WHERE id=$1;", res.GameID)
	}

	return nil
}

// func (rr *roundRepository) checkIfPlayerAlreadyPlayed(ctx context.Context, playerId int, roundId int) bool {
// 	var roundQuery string
// }

func (rr *roundRepository) UpdateHand(ctx context.Context, res *domain.RoundContext) error {

	player_query := `
		SELECT player_one_id, player_two_id FROM games WHERE id = $1;
	`
	type playerQueryRes struct {
		player_one int
		player_two int
	}
	var player_query_res playerQueryRes
	// Query Game player ids
	err := rr.db.QueryRowContext(ctx, player_query, &res.GameID).Scan(&player_query_res.player_one, &player_query_res.player_two)

	var set_player_hand_query string
	switch &res.CurrentPlayer {
	case &res.PlayerOneID:
		set_player_hand_query = `UPDATE rounds SET player_one_hand = $1 WHERE id = $2`
	case &res.PlayerTwoID:
		set_player_hand_query = "UPDATE rounds SET player_two_hand = $1 WHERE id = $2"
	default:
		return errors.New("Player does not belong here or is missing")
	}
	// Query to set the player hand
	_ = rr.db.QueryRowContext(ctx, set_player_hand_query, res.CurrentPlayerHand(), &res.ID)
	err = rr.CheckForWinner(ctx, res)
	if err != nil {
		return err
	}
	fmt.Println("CheckForWinner Response: ", res)
	return nil
}

/**
UTIL Functions
**/
// resolve winner from hands
func resolveHands(hand_one string, hand_two string) int {
	switch hand_one {
	case "rock":
		if hand_two == "paper" {
			return 2
		}
		if hand_two == "scissors" {
			return 1
		}
		if hand_two == "rock" {
			return 0
		}
	case "scissors":
		if hand_two == "rock" {
			return 2
		}
		if hand_two == "paper" {
			return 1
		}
		if hand_two == "scissors" {
			return 0
		}
	case "paper":
		if hand_two == "scissors" {
			return 2
		}
		if hand_two == "rock" {
			return 1
		}
		if hand_two == "paper" {
			return 0
		}
	}
	return 0
}
