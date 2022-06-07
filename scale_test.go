package jitter

import (
	"fmt"
	"testing"
	"time"
)

func TestScale(t *testing.T) {
	t.Run("sample distribution", func(t *testing.T) {
		const (
			d       = 10 * time.Second
			f       = 0.5
			samples = 20
		)

		for i := 0; i < samples; i++ {
			r := Scale(d, f)
			t.Log(r)
			if r < 5*time.Second || r > 15*time.Second {
				t.Error("sample outside of range: ", r)
			}
		}
	})

	t.Run("nonpositive scaling factors panic", func(t *testing.T) {
		var testcases = []struct {
			d         time.Duration
			f         float64
			wantPanic bool
		}{
			{d: time.Second, f: 0.1, wantPanic: false},
			{d: time.Second, f: 100.0, wantPanic: false},
			{d: time.Second, f: -0.1, wantPanic: true},
			{d: time.Second, f: 0.0, wantPanic: true},
		}

		for _, tc := range testcases {
			t.Run(fmt.Sprintf("%f", tc.f), func(t *testing.T) {
				assertPanic(t, func() { Scale(tc.d, tc.f) }, tc.wantPanic)
			})
		}
	})
}

func assertPanic(t *testing.T, f func(), want bool) {
	t.Helper()
	defer func() {
		r := recover()
		if got := (r != nil); got != want {
			t.Errorf("wantPanic: %v, gotPanic: %v", want, got)
		}
	}()
	f()
}

func BenchmarkScale(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scale(time.Second, 0.5)
	}
}
