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

func (ur *UserRepositoryImpl) Get(id int) (models.User, error) {
	row := ur.db.QueryRow(
		context.Background(),
		`SELECT id, email, username, password_hash, created_at, updated_at
		 FROM users 
		 WHERE id = $1;`,
		id,
	)
	return ur.fromRow(row)
}

func (ur *UserRepositoryImpl) Create(user models.User) (models.User, error) {
	row := ur.db.QueryRow(
		context.Background(),
		`INSERT INTO users (email, username, password_hash) 
		 VALUES ($1, $2, crypt($3, gen_salt('bf'))) 
    	 RETURNING id, email, username, password_hash, created_at, updated_at;`,
		user.Email, user.Username, user.Password,
	)
	return ur.fromRow(row)
}

func (ur *UserRepositoryImpl) fromRow(row pgx.Row) (models.User, error) {
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
