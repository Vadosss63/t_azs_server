package repository

import (
	"github.com/Vadosss63/t-azs/internal/repository/trbl_button"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/Vadosss63/t-azs/internal/repository/ya_pay"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool           *pgxpool.Pool
	YaPayRepo      *ya_pay.YaPayRepo
	UserRepo       *user.UserRepo
	TrblButtonRepo *trbl_button.TrblButtonRepo
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:           pool,
		YaPayRepo:      ya_pay.NewYaPayRepository(pool),
		UserRepo:       user.NewRepository(pool),
		TrblButtonRepo: trbl_button.NewRepository(pool),
	}
}
