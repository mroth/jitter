// Package jitter provides functionality for generating durations and tickers
// that deviate from true periodicity within specified bounds.
package jitter

/*
Note that if you are looking to use this to implement jitter for a backoff timer
for task retries, you may wish to check out https://github.com/kamilsk/retry
instead, which is very full featured and contains its own jitter support.
*/

import (
	"errors"
	"math"
	"time"
)

// Scale simulates jitter by scaling a time.Duration randomly within factor f.
//
// The duration d must be greater than zero; and the scaling factor f must be
// within the range 0 < f <= 1.0, or Scale will panic.
func Scale(d time.Duration, f float64) time.Duration {
	assertScaleArgs(d, f)
	min, max := scaleBounds(int64(d), f)
	return time.Duration(randRange(min, max))
}

func assertScaleArgs(d time.Duration, f float64) {
	switch {
	case d <= 0:
		panic(errors.New("non-positive interval for duration"))
	case f > 1.0 || f <= 0:
		panic(errors.New("scaling factor must be 0 < f <= 1.0"))
	}
}

// scaleBounds returns the min and max values for n after applying scaling
// factor f.
//
// if the max were to overflow, we instead truncate and return math.MaxInt64.
//
// as an internal function, it assumes n and f have already been validated via
// assertScaleArgs and does not handle edge cases outside of those parameters.
func scaleBounds(n int64, f float64) (min, max int64) {
	minf := math.Floor(float64(n) * (1 - f))
	maxf := math.Ceil(float64(n) * (1 + f))

	if maxf > math.MaxInt64 {
		return int64(minf), math.MaxInt64
	}
	return int64(minf), int64(maxf)
}
