package jitter

import (
	"context"
	"math"
	"testing"
	"testing/synctest"
	"time"
)

func TestNewTicker(t *testing.T) {
	type args struct {
		d      time.Duration
		factor float64
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "negative duration panics",
			args:      args{d: -1 * time.Second, factor: 0.1},
			wantPanic: true,
		},
		{
			name:      "negative factor panics",
			args:      args{d: time.Second, factor: -0.1},
			wantPanic: true,
		},
		{
			name:      "factor >1.0 panics",
			args:      args{d: time.Second, factor: 1.1},
			wantPanic: true,
		},
		{
			name:      "successful initiation",
			args:      args{d: time.Second, factor: 0.9},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewTicker did not panic")
					}
				}()
			}
			ticker := NewTicker(tt.args.d, tt.args.factor)
			ticker.Stop()
		})
	}
}

// checks time elapsed for sample number of ticks is within expected range
func TestTicker_start(t *testing.T) {
	const (
		d       = 10 * time.Millisecond
		factor  = 0.3
		samples = 10
	)

	synctest.Test(t, func(t *testing.T) {
		ticker := NewTicker(d, factor)
		t.Cleanup(ticker.Stop)

		t1 := time.Now()
		for range samples {
			<-ticker.C
		}

		var (
			elapsed = time.Since(t1)
			min     = time.Duration(math.Floor(float64(d)*(1-factor))) * samples
			max     = time.Duration(math.Ceil(float64(d)*(1+factor))) * samples
		)
		if elapsed < min || elapsed > max {
			t.Errorf("time elapsed for %v ticks %v outside of expected range %v - %v",
				samples, elapsed, min, max)
		}
	})
}

// checks that Stop() prevents further ticks from being sent
func TestTicker_stop(t *testing.T) {
	const (
		d            = time.Millisecond
		factor       = 0.1
		beforeTicks  = 3      // ticks before stop
		waitDuration = d * 10 // monitor after stop
	)

	synctest.Test(t, func(t *testing.T) {
		ticker := NewTicker(d, factor)
		for range beforeTicks {
			<-ticker.C
		}
		ticker.Stop()

		select {
		case <-ticker.C:
			t.Fatal("got tick after Stop()")
		case <-time.After(waitDuration):
			t.Log("detected no ticks after Stop()")
		}
	})
}

// checks that context cancellation prevents further ticks from being sent
func TestTicker_ctxExpired(t *testing.T) {
	const (
		d            = time.Millisecond
		factor       = 0.1
		beforeTicks  = 3      // ticks before cancel
		waitDuration = d * 10 // monitor after cancel
	)

	synctest.Test(t, func(t *testing.T) {
		ctx, cancelFunc := context.WithCancel(t.Context())
		ticker := NewTickerWithContext(ctx, d, factor)

		for range beforeTicks {
			<-ticker.C
		}

		cancelFunc()
		select {
		case <-ticker.C:
			t.Fatal("got tick after context cancelled")
		case <-time.After(waitDuration):
			t.Log("detected no ticks after context cancelled")
		}
	})
}
