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
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Scale simulates jitter by scaling a time.Duration randomly within factor f.
//
// The duration d must be greater than zero; and the scaling factor f must be
// within the range 0 < f <= 1.0, or Scale will panic.
func Scale(d time.Duration, f float64) time.Duration {
	assertScaleArgs(d, f)
	var (
		min = int64(math.Floor(float64(d) * (1 - f)))
		max = int64(math.Ceil(float64(d) * (1 + f)))
	)
	if min > max {
		panic(fmt.Sprintf("min > max, %v > %v", min, max))
	}
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

// randRange returns a nonnegative pseudo-random number in the half open
// interval [min, max) from the default Source.
//
// It panics if max < min.
func randRange(min, max int64) int64 {
	if min == max {
		return min
	}
	return rand.Int63n(max-min) + min
}
