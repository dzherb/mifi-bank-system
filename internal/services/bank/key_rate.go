package bank

import (
	"time"

	"github.com/dzherb/mifi-bank-system/internal/pkg/cache"
	"github.com/dzherb/mifi-bank-system/pkg/cbr"
)

const cacheKey = "cbr_key_rate"
const cachePeriod = time.Hour

type keyRateProvider func() (cbr.KeyRate, error)

var provider keyRateProvider = cbr.CurrentKeyRate

func setKeyRateProvider(p keyRateProvider) {
	provider = p
}

func CBRKeyRate() (cbr.KeyRate, error) {
	c := cache.Active()

	if item := c.Get(cacheKey); item != nil {
		return item.Value().(cbr.KeyRate), nil
	}

	kr, err := provider()
	if err != nil {
		return kr, err
	}

	c.Set(cacheKey, kr, cachePeriod)

	return kr, nil
}

const Margin = 5.

func KeyRate() (cbr.KeyRate, error) {
	kr, err := CBRKeyRate()
	if err != nil {
		return cbr.KeyRate{}, err
	}

	kr.Val += Margin
	kr.Date = time.Now().UTC()

	return kr, nil
}
