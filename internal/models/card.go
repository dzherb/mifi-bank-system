package models

import "time"

type Card struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Number    int       `json:"number"`
	Expires   time.Time `json:"expires"`
	CVV       string    `json:"cvv"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
