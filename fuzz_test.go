//go:build go1.18
// +build go1.18

package jitter

import (
	"testing"
	"time"
)

func FuzzScale(f *testing.F) {
	f.Add(int64(time.Second), 0.1)

	f.Fuzz(func(t *testing.T, d int64, f float64) {
		_ = Scale(time.Duration(d), f)
	})
}
