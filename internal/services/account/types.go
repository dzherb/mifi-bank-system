package account

import (
	"github.com/dzherb/mifi-bank-system/internal/models"
	repo "github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/shopspring/decimal"
)

type Service interface {
	Withdraw(from int, amount decimal.Decimal) (models.Account, error)
	Deposit(to int, amount decimal.Decimal) (models.Account, error)
	Transfer(from int, to int, amount decimal.Decimal) error
}

func NewService() Service {
	return &ServiceImpl{db: storage.Conn()}
}

type ServiceImpl struct {
	db storage.Connection
	ar repo.AccountRepository
	tr repo.TransactionRepository
}
