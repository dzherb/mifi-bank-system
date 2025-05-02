package repo

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type TransactionRepository interface {
	Create(transaction models.Transaction) (models.Transaction, error)
}

func NewTransactionRepository() TransactionRepository {
	return TransactionRepositoryImpl{storage.Conn()}
}

func NewTransactionRepositoryFromTx(tx pgx.Tx) TransactionRepository {
	return TransactionRepositoryImpl{tx}
}

type TransactionRepositoryImpl struct {
	db storage.Connection
}

func (t TransactionRepositoryImpl) Create(
	transaction models.Transaction,
) (models.Transaction, error) {
	err := pgxscan.Get(
		context.Background(),
		t.db,
		&transaction,
		`INSERT INTO transactions 
		(sender_account_id, receiver_account_id, type, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id, 
			sender_account_id, 
			receiver_account_id, 
			type, 
			amount, 
			created_at, 
			updated_at;`,
		transaction.SenderAccountID,
		transaction.ReceiverAccountID,
		transaction.Type,
		transaction.Amount,
	)
	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}
