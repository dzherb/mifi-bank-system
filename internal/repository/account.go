package repo

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type AccountRepository interface {
	Get(id int) (models.Account, error)
	Create(account models.Account) (models.Account, error)
	Update(account models.Account) (models.Account, error)
}

type AccountRepositoryImpl struct {
	db storage.Connection
}

func NewAccountRepository() AccountRepository {
	return &AccountRepositoryImpl{db: storage.Conn()}
}

func NewAccountRepositoryWithTx(tx pgx.Tx) AccountRepository {
	return &AccountRepositoryImpl{db: tx}
}

func (ar *AccountRepositoryImpl) Get(id int) (models.Account, error) {
	account := models.Account{}
	err := pgxscan.Get(
		context.Background(),
		ar.db,
		&account,
		`SELECT id, user_id, balance, created_at, updated_at
		 FROM accounts
		 WHERE id = $1;`,
		id,
	)

	if err != nil {
		return models.Account{}, err
	}

	return account, err
}

func (ar *AccountRepositoryImpl) Create(
	account models.Account,
) (models.Account, error) {
	err := pgxscan.Get(
		context.Background(),
		ar.db,
		&account,
		`INSERT INTO accounts (user_id)
		 VALUES ($1)
    	 RETURNING id, user_id, balance, created_at, updated_at;`,
		account.UserID,
	)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (ar *AccountRepositoryImpl) Update(
	account models.Account,
) (models.Account, error) {
	err := pgxscan.Get(
		context.Background(),
		ar.db,
		&account,
		`UPDATE accounts
		SET balance = $2
		WHERE id = $1
		RETURNING id, user_id, balance, created_at, updated_at;`,
		account.ID, account.Balance,
	)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
