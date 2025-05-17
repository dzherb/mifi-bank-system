package cbr

import (
	"time"

	"github.com/dzherb/mifi-bank-system/internal/pkg/cache"
	cbrpkg "github.com/dzherb/mifi-bank-system/pkg/cbr"
)

const cacheKey = "cbr_key_rate"
const cachePeriod = time.Hour

type keyRateProvider func() (cbrpkg.KeyRate, error)

var provider keyRateProvider = cbrpkg.CurrentKeyRate

func CurrentKeyRate() (cbrpkg.KeyRate, error) {
	c := cache.Active()

	if item := c.Get(cacheKey); item != nil {
		return item.Value().(cbrpkg.KeyRate), nil
	}

	kr, err := provider()
	if err != nil {
		return kr, err
	}

	c.Set(cacheKey, kr, cachePeriod)

	return kr, nil
}

func setKeyRateProvider(p keyRateProvider) {
	provider = p
}
