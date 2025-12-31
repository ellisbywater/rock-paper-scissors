package service

import (
	"context"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
)

type GameService struct {
	repo domain.GameRepository
}

func NewGameService(repo domain.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (gs *GameService) NewGame(ctx context.Context, total_rounds int, player_one_id int, player_two_id int) (*domain.GameCreateResponse, error) {
	game_req := domain.GameCreateRequest{
		TotalRounds: total_rounds,
		PlayerOneID: player_one_id,
		PlayerTwoID: player_two_id,
	}
	var game_res domain.GameCreateResponse
	err := gs.repo.Create(ctx, game_req, &game_res)
	if err != nil {
		return &game_res, err
	}
	return &game_res, nil
}

func (gs *GameService) GetGame(ctx context.Context, id int) (*domain.GameResponse, error) {
	var game domain.GameResponse
	err := gs.repo.Get(ctx, id, &game)
	if err != nil {
		return &game, err
	}
	return &game, nil
}

type PlayerService struct {
	repo domain.PlayerRepository
}

func NewPlayerRepository(repo domain.PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

func (ps *PlayerService) CreatePlayer(ctx context.Context, username string) (*domain.PlayerResponse, error) {
	player_req := domain.PlayerCreateRequest{
		UserName: username,
	}
	var player domain.PlayerResponse
	err := ps.repo.Create(ctx, player_req, &player)
	if err != nil {
		return &player, err
	}
	return &player, nil
}

func (ps *PlayerService) GetPlayer(ctx context.Context, id int) (*domain.PlayerResponse, error) {
	var player domain.PlayerResponse
	err := ps.repo.Get(ctx, id, &player)
	if err != nil {
		return &player, err
	}
	return &player, nil
}

func (ps *PlayerService) GetPlayerGames(ctx context.Context, id int) (*[]domain.GameResponse, error) {
	var games []domain.GameResponse
	err := ps.repo.GetGames(ctx, id, &games)
	if err != nil {
		return &games, err
	}
	return &games, nil
}

type RoundService struct {
	repo domain.RoundRepository
}

func NewRoundRepository(repo domain.RoundRepository) *RoundService {
	return &RoundService{repo: repo}
}

func (rs *RoundService) Create(ctx context.Context, req domain.RoundCreateRequest) (*domain.RoundCreateResponse, error) {
	var round_res domain.RoundCreateResponse
	err := rs.repo.Create(ctx, req, &round_res)
	if err != nil {
		return &round_res, err
	}
	return &round_res, nil
}

func (rs *RoundService) Get(ctx context.Context, id int) (*domain.Round, error) {
	var round_res domain.Round
	err := rs.repo.Get(ctx, id, &round_res)
	if err != nil {
		return &round_res, err
	}
	return &round_res, nil
}

func (rs *RoundService) UpdateHand(ctx context.Context, hand domain.RoundPlayerInput) (*domain.Round, error) {
	var round_res domain.Round
	err := rs.repo.UpdateHand(ctx, hand, &round_res)
	if err != nil {
		return &round_res, err
	}
	return &round_res, nil
}
