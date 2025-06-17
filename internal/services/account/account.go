package account

import (
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(userID int) (models.Account, error)
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

func (s *ServiceImpl) Create(userID int) (models.Account, error) {
	account := models.Account{UserID: userID}
	return repo.NewAccountRepository().Create(account)
}
