package account

import "errors"

var ErrNotEnoughMoney = errors.New("not enough money")

var ErrSameAccount = errors.New(
	"source and destination accounts must be different",
)

var ErrNotPositiveAmount = errors.New("amount must be positive")
