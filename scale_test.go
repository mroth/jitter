package jitter

import (
	"fmt"
	"math"
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

	t.Run("invalid arguments panic", func(t *testing.T) {
		var testcases = []struct {
			d         time.Duration
			f         float64
			wantPanic bool
		}{
			{d: time.Second, f: 0.1, wantPanic: false},
			{d: time.Second, f: 1.1, wantPanic: true},
			{d: time.Second, f: -0.1, wantPanic: true},
			{d: time.Second, f: 0.0, wantPanic: true},
			{d: -1 * time.Second, f: 0.1, wantPanic: true},              // negative duration
			{d: time.Duration(math.MaxInt64), f: 0.1, wantPanic: false}, // max duration
		}

		for _, tc := range testcases {
			t.Run(fmt.Sprintf("d=%d,f=%v", tc.d, tc.f), func(t *testing.T) {
				assertPanic(t, func() { Scale(tc.d, tc.f) }, tc.wantPanic)
			})
		}
	})
}

func Test_scaleBounds(t *testing.T) {
	type args struct {
		n int64
		f float64
	}
	tests := []struct {
		name    string
		args    args
		wantMin int64
		wantMax int64
	}{
		{name: "typical1", args: args{n: 1, f: 1}, wantMin: 0, wantMax: 2},
		{name: "typical2", args: args{n: 100, f: 0.5}, wantMin: 50, wantMax: 150},
		{name: "overflow truncate", args: args{n: math.MaxInt64, f: 0.5}, wantMin: math.MaxInt64/2 + 1, wantMax: math.MaxInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := scaleBounds(tt.args.n, tt.args.f)
			if (gotMin != tt.wantMin) || (gotMax != tt.wantMax) {
				t.Errorf("scaleBounds() = got (%d, %d), want (%d, %d)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func assertPanic(t *testing.T, f func(), want bool) {
	t.Helper()
	defer func() {
		r := recover()
		if got := (r != nil); got != want {
			t.Errorf("wantPanic: %v, gotPanic: %v [%v]", want, got, r)
		}
	}()
	f()
}

func BenchmarkScale(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scale(time.Second, 0.5)
	}
}
