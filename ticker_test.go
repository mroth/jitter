package jitter

import (
	"context"
	"math"
	"testing"
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

func TestTicker_start(t *testing.T) {
	t.Parallel()

	// measures actual timing, need to utilize a high enough time period to not
	// have scheduler overhead be a factor
	const (
		d        = 10 * time.Millisecond
		factor   = 0.3
		samples  = 10
		overhead = 1 * time.Millisecond
	)

	// check time elapsed for sample number of ticks is within expected range
	ticker := NewTicker(d, factor)
	t1 := time.Now()
	for range samples {
		<-ticker.C
	}

	var (
		elapsed = time.Since(t1)
		min     = time.Duration(math.Floor(float64(d)*(1-factor)))*samples - overhead
		max     = time.Duration(math.Ceil(float64(d)*(1+factor)))*samples + overhead
	)
	if elapsed < min || elapsed > max {
		t.Errorf("time elapsed for %v ticks %v outside of expected range %v - %v",
			samples, elapsed, min, max)
	}
}

func TestTicker_stop(t *testing.T) {
	t.Parallel()

	const (
		d            = time.Millisecond
		factor       = 0.1
		beforeTicks  = 3      // ticks before stop
		waitDuration = d * 10 // monitor after stop
	)

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
}

func TestTicker_ctxExpired(t *testing.T) {
	t.Parallel()

	const (
		d            = time.Millisecond
		factor       = 0.1
		beforeTicks  = 3      // ticks before cancel
		waitDuration = d * 10 // monitor after cancel
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
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
}
