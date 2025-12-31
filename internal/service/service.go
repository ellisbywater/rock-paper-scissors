package service

import "github.com/ellisbywater/http-rock-paper-scissors/internal/domain"

type gameService struct {
	repo domain.GameRepository
}
