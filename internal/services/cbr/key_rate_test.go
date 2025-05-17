package cbr_test

import (
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/pkg/cache"
	"github.com/dzherb/mifi-bank-system/internal/services/cbr"
	cbrpkg "github.com/dzherb/mifi-bank-system/pkg/cbr"
)

func TestCurrentKeyRate(t *testing.T) {
	cache.Init()
	t.Cleanup(cache.Close)

	const keyRate = 21

	cbr.SetKeyRateProvider(func() (cbrpkg.KeyRate, error) {
		return cbrpkg.KeyRate{Val: keyRate, Date: time.Now()}, nil
	})

	kr, err := cbr.CurrentKeyRate()
	if err != nil {
		t.Error(err)
		return
	}

	if kr.Val != keyRate {
		t.Errorf("expected key rate to be %v, got %v", keyRate, kr.Val)
	}

	// Check that the key rate value is cached
	cbr.SetKeyRateProvider(func() (cbrpkg.KeyRate, error) {
		return cbrpkg.KeyRate{Val: 25, Date: time.Now()}, nil
	})

	kr, err = cbr.CurrentKeyRate()
	if err != nil {
		t.Error(err)
		return
	}

	if kr.Val != keyRate {
		t.Errorf("expected key rate to be %v, got %v", keyRate, kr.Val)
	}
}
