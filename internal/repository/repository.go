package repository

import (
	"database/sql"

	"github.com/ellisbywater/http-rock-paper-scissors/internal/domain"
)

type GameRepository struct {
	db *sql.DB
}

func (gr *GameRepository) Create(game domain.GameCreateRequest) (error, domain.GameResponse) {

}
