package repository

import (
	"context"
	"database/sql"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
)

type GameRepository struct {
	db *sql.DB
}

func (gr *GameRepository) Create(ctx context.Context, game domain.GameCreateRequest, res *domain.GameCreateResponse) error {
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

func (gr *GameRepository) Get(ctx context.Context, id int, res *domain.GameResponse) error {
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

type PlayerRepository struct {
	db *sql.DB
}
