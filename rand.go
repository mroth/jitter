package jitter

import "math/rand/v2"

// randRange returns a nonnegative pseudo-random number in the half open
// interval [min, max) from the default Source.
//
// It panics if max < min.
func randRange(min, max int64) int64 {
	if min == max {
		return min
	}
	return rand.Int64N(max-min) + min
}
