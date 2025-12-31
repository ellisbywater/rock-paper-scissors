package repository

import (
	"context"
	"database/sql"

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
			player_one,
			player_two
		) Values (
		 	$1,
			$2,
			$3
		 )
		RETURNING id, created_at, player_one, player_two;
	`
	err := gr.db.QueryRowContext(
		ctx,
		query,
		game.TotalRounds,
		game.PlayerOneID,
		game.PlayerTwoID,
	).Scan(&res.ID, &res.CreatedAt, &res.PlayerOneId, &game.PlayerTwoID)

	if err != nil {
		return err
	}
	return nil
}

func (gr *gameRepository) Get(ctx context.Context, id int, res *domain.GameResponse) error {
	query := `
		SELECT id, total_rounds, player_one, player_two, winner, created_at
		FROM games
		WHERE id = $1
	`
	err := gr.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.TotalRounds,
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
		) RETURNING id, username;
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
		SELECT * FROM games g WHERE player_one = $1 OR player_two = $1
		 JOIN rounds r ON g.id = r.game;
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
		SELECT * FROM rounds WHERE id=$1
	`
	err := rr.db.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.Count, &res.PlayerOneHand, &res.PlayerTwoHand)
	if err != nil {
		return err
	}
	return nil
}

func (rr *roundRepository) Create(ctx context.Context, round_create_request domain.RoundCreateRequest, res *domain.RoundCreateResponse) error {
	query := `
		INSERT INTO rounds (
			game,
			count,
			player_one,
			player_two
		) VALUES (
		 	$1,
			$2,
			$3,
			$4
		) RETURNING id;
	`
	err := rr.db.QueryRowContext(
		ctx,
		query,
		round_create_request.GameId,
		round_create_request.Count,
		round_create_request.PlayerOneId,
		round_create_request.PlayerTwoID,
	).Scan(&res.Id)

	if err != nil {
		return err
	}
	return nil
}

func (rr *roundRepository) UpdateHand(ctx context.Context, player_input domain.RoundPlayerInput, res *domain.Round) error {
	query := `
		UPDATE rounds SET $1 = $2 WHERE id = $3 RETURNING *;
	`
	var player string
	if player_input.IsPlayerOne {
		player = "player_one_hand"
	} else {
		player = "player_two_hand"
	}
	err := rr.db.QueryRowContext(ctx, query, player, player_input.Hand, player_input.RoundId).Scan(
		&res.ID,
		&res.Count,
		&res.PlayerOneHand,
		&res.PlayerTwoHand,
		&res.Winner,
	)
	if err != nil {
		return err
	}
	return nil
}
