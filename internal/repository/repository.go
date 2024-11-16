package repository

import (
	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/azs_button"
	"github.com/Vadosss63/t-azs/internal/repository/azs_statistics"
	"github.com/Vadosss63/t-azs/internal/repository/receipt"
	"github.com/Vadosss63/t-azs/internal/repository/trbl_button"
	"github.com/Vadosss63/t-azs/internal/repository/updater_button"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/Vadosss63/t-azs/internal/repository/ya_azs"
	"github.com/Vadosss63/t-azs/internal/repository/ya_pay"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	AzsRepo           azs.AzsRepository
	AzsButtonRepo     azs_button.AzsButtonRepository
	ReceiptRepo       receipt.ReceiptRepository
	AzsStatRepo       azs_statistics.StatisticsRepository
	TrblButtonRepo    trbl_button.TrblButtonRepository
	UpdaterButtonRepo updater_button.UpdaterButtonRepository
	UserRepo          user.UserRepository
	YaAzsRepo         ya_azs.YaAzsRepository
	YaPayRepo         ya_pay.YaPayRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		AzsRepo:           azs.NewRepository(pool),
		AzsButtonRepo:     azs_button.NewRepository(pool),
		ReceiptRepo:       receipt.NewRepository(pool),
		AzsStatRepo:       azs_statistics.NewStatisticsRepository(pool),
		TrblButtonRepo:    trbl_button.NewRepository(pool),
		UpdaterButtonRepo: updater_button.NewRepository(pool),
		UserRepo:          user.NewRepository(pool),
		YaAzsRepo:         ya_azs.NewRepository(pool),
		YaPayRepo:         ya_pay.NewRepository(pool),
	}
}
