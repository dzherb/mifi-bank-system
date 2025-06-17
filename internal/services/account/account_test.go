package account_test

import (
	"errors"
	"os"
	"testing"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/services/account"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.RunTestsWithTempDB(func() int {
		return storage.RunTestsWithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func prepareTransferAccounts() (int, int, error) {
	users := []models.User{
		{
			Email:    "test1@test.com",
			Username: "test1",
			Password: "test1",
		},
		{
			Email:    "test2@test.com",
			Username: "test2",
			Password: "test1",
		},
	}

	createdUsers := make([]models.User, 0, len(users))
	ur := repo.NewUserRepository()

	for _, user := range users {
		u, err := ur.Create(user)
		if err != nil {
			return 0, 0, err
		}

		createdUsers = append(createdUsers, u)
	}

	accounts := []models.Account{
		{
			UserID:  createdUsers[0].ID,
			Balance: decimal.NewFromInt(500),
		},
		{
			UserID:  createdUsers[1].ID,
			Balance: decimal.NewFromInt(1000),
		},
	}

	createdAccounts := make([]models.Account, 0, len(accounts))
	ar := repo.NewAccountRepository()

	for _, a := range accounts {
		ac, err := ar.Create(a)
		if err != nil {
			return 0, 0, err
		}

		createdAccounts = append(createdAccounts, ac)
	}

	return createdAccounts[0].ID, createdAccounts[1].ID, nil
}

func TestServiceImpl_Transfer(t *testing.T) {
	storage.TestWithTransaction(t)

	a1, a2, err := prepareTransferAccounts()
	if err != nil {
		t.Fatal(err)
	}

	as := account.NewService()

	err = as.Transfer(a1, a2, decimal.NewFromInt(100))
	if err != nil {
		t.Errorf("failed to transfer: %v", err)
		return
	}

	ar := repo.NewAccountRepository()

	a1Refreshed, err := ar.Get(a1)
	if err != nil {
		t.Fatal(err)
	}

	a2Refreshed, err := ar.Get(a2)
	if err != nil {
		t.Fatal(err)
	}

	expected := decimal.NewFromInt(400)
	if !a1Refreshed.Balance.Equal(expected) {
		t.Errorf("expected balance %v, got %v", expected, a1Refreshed.Balance)
	}

	expected = decimal.NewFromInt(1100)
	if !a2Refreshed.Balance.Equal(expected) {
		t.Errorf("expected balance %v, got %v", expected, a2Refreshed.Balance)
	}
}

func TestServiceImpl_TransferFail(t *testing.T) {
	storage.TestWithTransaction(t)

	accountID, _, err := prepareTransferAccounts()
	if err != nil {
		t.Fatal(err)
	}

	invalidID := 42

	as := account.NewService()

	err = as.Transfer(accountID, invalidID, decimal.NewFromInt(100))
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}

func TestServiceImplBalance(t *testing.T) { //nolint:gocognit
	cases := []struct {
		amount   decimal.Decimal
		expected decimal.Decimal
		err      error
		as       models.TransactionType
	}{
		{
			amount:   decimal.NewFromInt(100),
			expected: decimal.NewFromInt(600),
			err:      nil,
			as:       models.Deposit,
		},
		{
			amount:   decimal.NewFromInt(-100),
			expected: decimal.NewFromInt(0),
			err:      account.ErrNotPositiveAmount,
			as:       models.Deposit,
		},
		{
			amount:   decimal.NewFromInt(0),
			expected: decimal.NewFromInt(0),
			err:      account.ErrNotPositiveAmount,
			as:       models.Deposit,
		},
		{
			amount:   decimal.NewFromInt(100),
			expected: decimal.NewFromInt(400),
			err:      nil,
			as:       models.Withdrawal,
		},
		{
			amount:   decimal.NewFromInt(500),
			expected: decimal.NewFromInt(0),
			err:      nil,
			as:       models.Withdrawal,
		},
		{
			amount:   decimal.NewFromInt(501),
			expected: decimal.NewFromInt(0),
			err:      account.ErrNotEnoughMoney,
			as:       models.Withdrawal,
		},
		{
			amount:   decimal.NewFromInt(-100),
			expected: decimal.NewFromInt(0),
			err:      account.ErrNotPositiveAmount,
			as:       models.Withdrawal,
		},
		{
			amount:   decimal.NewFromInt(0),
			expected: decimal.NewFromInt(0),
			err:      account.ErrNotPositiveAmount,
			as:       models.Withdrawal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.amount.String(), func(t *testing.T) {
			storage.TestWithTransaction(t)

			accountID, _, err := prepareTransferAccounts()
			if err != nil {
				t.Fatal(err)
			}

			as := account.NewService()

			var accountUpdated models.Account

			if tc.as == models.Deposit {
				accountUpdated, err = as.Deposit(accountID, tc.amount)
			} else {
				accountUpdated, err = as.Withdraw(accountID, tc.amount)
			}

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}

				return
			}

			if tc.err == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if accountUpdated.ID != accountID {
				t.Errorf("expected ID %v, got %v", accountID, accountUpdated.ID)
			}

			if !accountUpdated.Balance.Equal(tc.expected) {
				t.Errorf(
					"expected balance %v, got %v",
					tc.expected,
					accountUpdated.Balance,
				)
			}
		})
	}
}

func TestServiceImplCreate(t *testing.T) {
	ur := repo.NewUserRepository()
	user, err := ur.Create(models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "123456",
	})

	if err != nil {
		t.Fatal(err)
	}

	service := account.NewService()

	acc, err := service.Create(user.ID)
	if err != nil {
		t.Errorf("failed to create account: %v", err)
	}

	if acc.UserID != user.ID {
		t.Errorf("expected user ID %v, got %v", user.ID, acc.UserID)
	}

	if !acc.Balance.IsZero() {
		t.Errorf("expected balance %v, got %v", 0, acc.Balance)
	}
}
