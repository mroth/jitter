//go:build go1.18
// +build go1.18

package jitter

import (
	"math"
	"testing"
	"time"
)

func FuzzScale(f *testing.F) {
	f.Add(int64(time.Second), 0.1)
	f.Add(int64(1), 0.1)
	f.Add(int64(1), 1.0)
	f.Add(int64(math.MaxInt64), 0.1)
	f.Add(int64(math.MaxInt64), 1.0)

	f.Fuzz(func(t *testing.T, d int64, f float64) {
		if f <= 0 || f > 1.0 || d <= 0 {
			t.Skip() // documented panics
		}
		_ = Scale(time.Duration(d), f)
	})
}
