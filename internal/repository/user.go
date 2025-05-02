package repo

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Get(id int) (models.User, error)
	GetByCredentials(email, password string) (models.User, error)
	Create(user models.User) (models.User, error)
}

type UserRepositoryImpl struct {
	db storage.Connection
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		db: storage.Conn(),
	}
}

func NewUserRepositoryFromTx(tx pgx.Tx) UserRepository {
	return &UserRepositoryImpl{
		db: tx,
	}
}

func (ur *UserRepositoryImpl) Get(id int) (models.User, error) {
	user := models.User{}
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1;`,
		id,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) Create(user models.User) (models.User, error) {
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, crypt($3, gen_salt('bf')))
		RETURNING id, email, username, password_hash, created_at, updated_at;`,
		user.Email, user.Username, user.Password,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) GetByCredentials(
	email, password string,
) (models.User, error) {
	user := models.User{}
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1 AND crypt($2, password_hash) = password_hash;`,
		email, password,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
