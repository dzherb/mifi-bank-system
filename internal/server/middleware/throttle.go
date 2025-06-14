package middleware

import (
	"fmt"
	"net/http"

	"github.com/dzherb/mifi-bank-system/internal/server/responses"
	log "github.com/sirupsen/logrus"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

const (
	maxRatePerMinute = 100
	maxBurst         = 10
	lruCacheSize     = 1 << 13
)

func RateLimiter() func(http.Handler) http.Handler {
	store, err := memstore.NewCtx(lruCacheSize)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(maxRatePerMinute),
		MaxBurst: maxBurst,
	}

	rateLimiter, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	return (&throttled.HTTPRateLimiterCtx{
		DeniedHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			responses.WriteError(w, fmt.Errorf("limit exceeded"))
		}),
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}).RateLimit
}
