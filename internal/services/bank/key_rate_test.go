package bank_test

import (
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/cache"
	"github.com/dzherb/mifi-bank-system/internal/services/bank"
	"github.com/dzherb/mifi-bank-system/pkg/cbr"
)

func TestCBRKeyRate(t *testing.T) {
	cache.Init()
	t.Cleanup(cache.Close)

	const keyRate = 21

	bank.SetKeyRateProvider(func() (cbr.KeyRate, error) {
		return cbr.KeyRate{Val: keyRate, Date: time.Now()}, nil
	})

	kr, err := bank.CBRKeyRate()
	if err != nil {
		t.Error(err)
		return
	}

	if kr.Val != keyRate {
		t.Errorf("expected key rate to be %v, got %v", keyRate, kr.Val)
	}

	// Check that the key rate value is cached
	bank.SetKeyRateProvider(func() (cbr.KeyRate, error) {
		return cbr.KeyRate{Val: 25, Date: time.Now()}, nil
	})

	kr, err = bank.CBRKeyRate()
	if err != nil {
		t.Error(err)
		return
	}

	if kr.Val != keyRate {
		t.Errorf("expected key rate to be %v, got %v", keyRate, kr.Val)
	}
}

func TestKeyRate(t *testing.T) {
	cache.Init()
	t.Cleanup(cache.Close)

	const keyRate = 21

	bank.SetKeyRateProvider(func() (cbr.KeyRate, error) {
		return cbr.KeyRate{Val: keyRate, Date: time.Now()}, nil
	})

	kr, err := bank.KeyRate()
	if err != nil {
		t.Error(err)
		return
	}

	exp := keyRate + bank.Margin
	if kr.Val != exp {
		t.Errorf("expected key rate to be %v, got %v", exp, kr.Val)
	}
}
