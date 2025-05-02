package account

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func (s *ServiceImpl) Transfer(from int, to int, amount decimal.Decimal) error {
	err := storage.WithTransaction(
		context.Background(),
		s.db,
		func(tx pgx.Tx) error {
			s.prepareForTx(tx)
			return s.transferAndLogTransaction(from, to, amount)
		},
	)

	logFields := log.Fields{
		"from_account_id": from,
		"to_account_id":   to,
		"amount":          amount,
	}

	if err != nil {
		logFields["err"] = err
		log.WithFields(logFields).Error("Failed to transfer money")

		return err
	}

	log.WithFields(logFields).Info("Money successfully transferred")

	return nil
}

func (s *ServiceImpl) transferAndLogTransaction(
	from int,
	to int,
	amount decimal.Decimal,
) error {
	if from == to {
		return ErrSameAccount
	}

	_, err := s.withdraw(from, amount)
	if err != nil {
		return err
	}

	_, err = s.deposit(to, amount)
	if err != nil {
		return err
	}

	_, err = s.tr.Create(
		models.Transaction{
			SenderAccountID:   &from,
			ReceiverAccountID: &to,
			Type:              models.Transfer,
			Amount:            amount,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
