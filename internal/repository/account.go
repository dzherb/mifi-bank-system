package repo

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
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

func (ar *AccountRepositoryImpl) Get(id int) (models.Account, error) {
	row := ar.db.QueryRow(
		context.Background(),
		`SELECT id, user_id, balance, created_at, updated_at
		 FROM accounts
		 WHERE id = $1;`,
		id,
	)

	return ar.fromRow(row)
}

func (ar *AccountRepositoryImpl) Create(
	account models.Account,
) (models.Account, error) {
	row := ar.db.QueryRow(
		context.Background(),
		`INSERT INTO accounts (user_id)
		 VALUES ($1)
    	 RETURNING id, user_id, balance, created_at, updated_at;`,
		account.UserID,
	)

	return ar.fromRow(row)
}

func (ar *AccountRepositoryImpl) Update(
	account models.Account,
) (models.Account, error) {
	row := ar.db.QueryRow(
		context.Background(),
		`UPDATE accounts
		SET balance = $2
		WHERE id = $1
		RETURNING id, user_id, balance, created_at, updated_at;`,
		account.ID, account.Balance,
	)

	return ar.fromRow(row)
}

func (ar *AccountRepositoryImpl) fromRow(row pgx.Row) (models.Account, error) {
	account := models.Account{}
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
