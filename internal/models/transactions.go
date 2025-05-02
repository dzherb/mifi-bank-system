package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	Withdrawal TransactionType = "withdrawal"
	Deposit    TransactionType = "deposit"
	Transfer   TransactionType = "transfer"
)

type Transaction struct {
	ID                int             `json:"id"`
	SenderAccountID   *int            `json:"sender_account_id"`
	ReceiverAccountID *int            `json:"receiver_account_id"`
	Type              TransactionType `json:"type"`
	Amount            decimal.Decimal `json:"amount"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
