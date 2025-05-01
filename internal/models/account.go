package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
