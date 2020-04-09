// Package jitter provides functionality for generating durations and tickers
// that deviate from true periodicity within specified bounds.
//
// All functionality in this package depends on global rand, so you will want to
// seed it before utilization.
package jitter

/*
Note that if you are looking to use this to implement jitter for a backoff timer
for task retries, you may wish to check out https://github.com/kamilsk/retry
instead, which is very full featured and contains its own jitter support.
*/

import (
	"math"
	"math/rand"
	"time"
)

// Scale simulates jitter by scaling a time.Duration randomly within factor.
//
// Note that using a factor of `math.Abs(f) > 1.0` may result in the sign of the
// result changing (e.g. a positive Duration may become negative, and vice
// versa). Additionally, a factor of `math.Abs(f) == 1.0` may result in a zero
// Duration. If you wish to avoid these potential scenarios, confine your factor
// such that `0.0 < math.Abs(f) < 1.0`.
func Scale(d time.Duration, factor float64) time.Duration {
	min := int64(math.Floor(float64(d) * (1 - factor)))
	max := int64(math.Ceil(float64(d) * (1 + factor)))
	return time.Duration(randRange(min, max))
}

// canonical Go boilerplate for random within a range
func randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
