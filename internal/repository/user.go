package repository

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Get(id int) (models.User, error)
	Create(user models.User) (models.User, error)
}

type UserRepositoryImpl struct {
	db pgx.Tx
}

func NewUserRepository(tx pgx.Tx) UserRepository {
	return &UserRepositoryImpl{
		db: tx,
	}
}

func (u *UserRepositoryImpl) Get(id int) (models.User, error) {
	row := u.db.QueryRow(
		context.Background(),
		`SELECT id, email, username, password, created_at, updated_at
		 FROM users 
		 WHERE id = $1;`,
		id,
	)
	return u.fromRow(row)
}

func (u *UserRepositoryImpl) Create(user models.User) (models.User, error) {
	row := u.db.QueryRow(
		context.Background(),
		`INSERT INTO users (email, username, password) 
		 VALUES ($1, $2, $3) 
    	 RETURNING id, email, username, password, created_at, updated_at;`,
		user.Email, user.Username, user.Password,
	)
	return u.fromRow(row)
}

func (u *UserRepositoryImpl) fromRow(row pgx.Row) (models.User, error) {
	user := models.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
