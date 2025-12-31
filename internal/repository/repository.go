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

func (rr *roundRepository) Get(ctx context.Context, id int, res *domain.Round) error {
	query := `
		SELECT * FROM rounds WHERE id=$1;
	`
	err := rr.db.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.Count, &res.PlayerOneHand, &res.PlayerTwoHand)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Check if game is at limit of total rounds
func (rr *roundRepository) Create(ctx context.Context, round_create_request domain.RoundCreateRequest, res *domain.RoundCreateResponse) error {
	type checkResult struct {
		current_round int
		total_rounds  int
		finished      bool
		player_one_id int
		player_two_id int
	}
	var checkCountResult checkResult
	check_count_query := `
		SELECT current_round, total_rounds, player_one_id, player_two_id, finished FROM games WHERE id=$1;
	`
	err := rr.db.QueryRowContext(ctx, check_count_query, round_create_request.GameId).Scan(&checkCountResult.current_round, &checkCountResult.total_rounds, &checkCountResult.player_one_id, &checkCountResult.player_two_id, &checkCountResult.finished)
	if err != nil {
		return err
	}
	if checkCountResult.finished {
		return errors.New("Game Already Finished!")
	}
	if checkCountResult.current_round == checkCountResult.total_rounds {
		return fmt.Errorf("Already on last round: %d", checkCountResult.current_round)
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
		) RETURNING id;
	`
	err = rr.db.QueryRowContext(
		ctx,
		query,
		round_create_request.GameId,
		checkCountResult.current_round+1,
		checkCountResult.player_one_id,
		checkCountResult.player_two_id,
	).Scan(&res.Id)

	if err != nil {
		return err
	}
	return nil
}

// TODO: Has to be someway to simplify or condense the queries
// Checks For Winner
// Updates Score
// Updates game finished
func (rr *roundRepository) CheckForWinner(ctx context.Context, round_id int, game_id int, res *domain.Round) error {
	type handsResults struct {
		player_one_id   int
		player_two_id   int
		player_one_hand domain.Hand
		player_two_hand domain.Hand
		message         string
	}
	var results handsResults
	query_player_select := `
		SELECT player_one_id, player_two_id, player_one_hand, player_two_hand FROM round WHERE id = $1 RETURNING *;
	`
	// Retrieve Fields for comparison
	err := rr.db.QueryRowContext(ctx, query_player_select, round_id).Scan(&results.player_one_id, &results.player_two_id, &results.player_one_hand, &results.player_two_hand)
	if err != nil {
		return err
	}

	var winner_query string      //dynamic query string
	var player_score_name string //dynamic query string
	// Calculate score and create respective queries
	if len(results.player_one_hand) != 0 && len(results.player_two_hand) != 0 {
		winner_num := resolveHands(string(results.player_one_hand), string(results.player_two_hand))
		switch winner_num {
		case 1:
			winner_query = formatWinnerQuery(&results.player_one_id, round_id)
			player_score_name = "player_one_score"
		case 2:
			winner_query = formatWinnerQuery(&results.player_two_id, round_id)
			player_score_name = "player_two_score"
		default:
			winner_query = fmt.Sprintf("UPDATE rounds SET winner=NULL, finished=True WHERE id=%d RETURNING *;", round_id)
			player_score_name = ""
		}
	} else {
		return nil
	}

	type checkGameCount struct {
		total_rounds  int
		current_round int
	}
	var gameCount checkGameCount
	// Update current round
	err = rr.db.QueryRowContext(ctx, "UPDATE games SET current_round = current_round + 1 WHERE id=$1 RETURNING current_round, total_rounds;", game_id).Scan(&gameCount.current_round, &gameCount.total_rounds)
	if err != nil {
		return err
	}
	// Update Round Winner
	err = rr.db.QueryRowContext(ctx, winner_query).Scan(&res.ID, &res.GameId, &res.Count, &res.PlayerOneHand, &res.PlayerTwoHand, &res.Winner, &res.Finished)
	if err != nil {
		return err
	}

	// Update Score
	if len(player_score_name) != 0 {
		score_query := formatUpdateGameScore(player_score_name, game_id)
		_ = rr.db.QueryRowContext(ctx, score_query)
	}

	if gameCount.current_round == gameCount.total_rounds {
		_ = rr.db.QueryRowContext(ctx, "UPDATE games SET finished=True WHERE id=$1;", game_id)
	}

	return nil
}

func (rr *roundRepository) UpdateHand(ctx context.Context, player_input domain.RoundPlayerInput, res *domain.Round) error {

	var player string
	player_query := `
		SELECT player_one_id, player_two_id FROM games WHERE id = $1;
	`
	type playerQueryRes struct {
		player_one int
		player_two int
	}
	var player_query_res playerQueryRes
	err := rr.db.QueryRowContext(ctx, player_query, player_input.GameID).Scan(&player_query_res.player_one, &player_query_res.player_two)

	switch player_input.PlayerID {
	case player_query_res.player_one:
		player = "player_one_hand"
	case player_query_res.player_two:
		player = "player_two_hand"
	default:
		return errors.New("Player does not belong here")
	}

	query := `
		UPDATE rounds SET $1 = $2 WHERE id = $3 RETURNING *;
	`

	_ = rr.db.QueryRowContext(ctx, query, player, player_input.Hand, player_input.RoundId)

	err = rr.CheckForWinner(ctx, player_input.RoundId, player_input.GameID, res)

	if err != nil {
		return err
	}
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

// format winner update query dynamically in a switch statement
func formatWinnerQuery(player_id *int, round_id int) string {
	return fmt.Sprintf("UPDATE rounds SET winner = %d, finished = True WHERE id = %d;", player_id, round_id)
}

// format game score update query dynamically in a switch statement
func formatUpdateGameScore(player_score_name string, game_id int) string {
	return fmt.Sprintf("UPDATE games SET %s = %s + 1 WHERE id=%d;", player_score_name, player_score_name, game_id)
}
