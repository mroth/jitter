package jitter

import (
	"context"
	"errors"
	"time"
)

// A Ticker holds a channel that delivers `ticks' of a clock at intervals,
// deviating with constrained random jitter.
//
// It adjusts the intervals or drops ticks to make up for slow receivers.
type Ticker struct {
	C      <-chan time.Time // The channel on which the ticks are delivered.
	cancel context.CancelFunc
}

// NewTicker returns a new Ticker containing a channel that will send the time
// with a period specified by the duration argument, but adjusted with random
// jitter based on the specified scaling factor.
//
// The duration d must be greater than zero; and the scaling factor f must be
// within the range 0 < f <= 1.0, or NewTicker will panic.
//
// Stop the ticker to release associated resources.
func NewTicker(d time.Duration, f float64) *Ticker {
	return NewTickerWithContext(context.Background(), d, f)
}

// NewTickerWithContext is identical to NewTicker but also takes a specified context.
// If this context is cancelled, the Ticker will automatically Stop.
func NewTickerWithContext(ctx context.Context, d time.Duration, f float64) *Ticker {
	switch {
	case d <= 0:
		panic(errors.New("non-positive interval for duration"))
	case f > 1.0 || f <= 0:
		panic(errors.New("factor must be 0 < f <= 1.0"))
	}

	// Add internal cancelFunc to the context, to be stored for use in Stop().
	ctx, cf := context.WithCancel(ctx)

	// Give the channel a 1-element time buffer.
	//
	// Similar to time.Timer, if the client falls behind while reading, we will
	// need to drop ticks on the floor until the client catches up.
	c := make(chan time.Time, 1)

	go func() {
		timer := time.NewTimer(Scale(d, f)) // initial timer

		for {
			select {
			case tc := <-timer.C:
				// Reset the internal timer for the next duration.
				//
				// Since the program has already received a value from timer.C,
				// the timer is known to have expired and the channel drained,
				// so t.Reset can be used directly.
				timer.Reset(Scale(d, f))

				// Non-blocking rebroadcast of time on c.
				//
				// Dropping sends on the floor is the desired behavior when the
				// reader gets behind, because the sends are periodic.
				//
				// c.f. func sendTime in https://golang.org/src/time/sleep.go.
				select {
				case c <- tc:
				default:
				}
			case <-ctx.Done():
				// stop the internal timer and drain its channel if needed
				if !timer.Stop() {
					<-timer.C
				}
				// ...then exit the goroutine
				return
			}
		}
	}()

	return &Ticker{C: c, cancel: cf}
}

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a concurrent goroutine reading from the
// channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	t.cancel()
}
