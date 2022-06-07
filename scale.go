// Package jitter provides functionality for generating durations and tickers
// that deviate from true periodicity within specified bounds.
//
// All functionality in this package currently utilizes global rand, so you will
// want to seed it before utilization.
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

// Scale simulates jitter by scaling a time.Duration randomly within factor f.
//
// Note that using a factor of f > 1.0 may result in the sign of the
// result changing (e.g. a positive Duration may become negative, and vice
// versa). Additionally, a factor of f == 1.0 may result in a zero
// Duration. If you wish to avoid these potential scenarios, confine your factor
// such that 0.0 < f < 1.0.
//
// If f <= 0, Scale will panic.
func Scale(d time.Duration, f float64) time.Duration {
	if f <= 0 {
		panic("invalid scaling factor for Scale")
	}

	var (
		min = int64(math.Floor(float64(d) * (1 - f)))
		max = int64(math.Ceil(float64(d) * (1 + f)))
	)
	return time.Duration(randRange(min, max))
}

// randRange is canonical Go boilerplate for random within a range.
//
// panics if max <= min.
func randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
