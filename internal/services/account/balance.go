package account

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	repo "github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func (s *ServiceImpl) Withdraw(
	from int,
	amount decimal.Decimal,
) (models.Account, error) {
	account := models.Account{}
	err := storage.WithTransaction(
		context.Background(),
		s.db,
		func(tx pgx.Tx) error {
			var err error

			s.prepareForTx(tx)
			account, err = s.withdrawAndLogTransaction(from, amount)

			return err
		},
	)

	logFields := log.Fields{
		"account_id": from,
		"amount":     amount,
	}

	if err != nil {
		logFields["err"] = err
		log.WithFields(logFields).Error("Withdraw failed")

		return account, err
	}

	log.WithFields(logFields).Info("Withdraw succeeded")

	return account, nil
}

func (s *ServiceImpl) Deposit(
	to int,
	amount decimal.Decimal,
) (models.Account, error) {
	account := models.Account{}
	err := storage.WithTransaction(
		context.Background(),
		s.db,
		func(tx pgx.Tx) error {
			var err error

			s.prepareForTx(tx)
			account, err = s.depositAndLogTransaction(to, amount)

			return err
		},
	)

	logFields := log.Fields{
		"account_id": to,
		"amount":     amount,
	}

	if err != nil {
		logFields["err"] = err
		log.WithFields(logFields).Error("Deposit failed")

		return account, err
	}

	log.WithFields(logFields).Info("Deposit succeeded")

	return account, nil
}

func (s *ServiceImpl) withdraw(
	from int,
	amount decimal.Decimal,
) (models.Account, error) {
	if !amount.IsPositive() {
		return models.Account{}, ErrNotPositiveAmount
	}

	account, err := s.ar.Get(from)
	if err != nil {
		return account, err
	}

	newBalance := account.Balance.Sub(amount)
	if newBalance.IsNegative() {
		return account, ErrNotEnoughMoney
	}

	account.Balance = newBalance

	return s.ar.Update(account)
}

func (s *ServiceImpl) withdrawAndLogTransaction(
	from int,
	amount decimal.Decimal,
) (models.Account, error) {
	account, err := s.withdraw(from, amount)
	if err != nil {
		return account, err
	}

	_, err = s.tr.Create(models.Transaction{
		ReceiverAccountID: &from,
		Type:              models.Withdrawal,
		Amount:            amount,
	})

	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *ServiceImpl) deposit(
	to int,
	amount decimal.Decimal,
) (models.Account, error) {
	if !amount.IsPositive() {
		return models.Account{}, ErrNotPositiveAmount
	}

	account, err := s.ar.Get(to)
	if err != nil {
		return account, err
	}

	account.Balance = account.Balance.Add(amount)

	return s.ar.Update(account)
}

func (s *ServiceImpl) depositAndLogTransaction(
	to int,
	amount decimal.Decimal,
) (models.Account, error) {
	account, err := s.deposit(to, amount)
	if err != nil {
		return account, err
	}

	_, err = s.tr.Create(models.Transaction{
		ReceiverAccountID: &to,
		Type:              models.Deposit,
		Amount:            amount,
	})

	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *ServiceImpl) prepareForTx(tx pgx.Tx) {
	s.ar = repo.NewAccountRepositoryFromTx(tx)
	s.tr = repo.NewTransactionRepositoryFromTx(tx)
}
