package security

import "time"

func timeToFloat64(t time.Time) float64 {
	return float64(t.UTC().UnixNano()) / 1e9
}
