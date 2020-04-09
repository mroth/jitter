package jitter

import (
	"context"
	"errors"
	"time"
)

// A Ticker holds a channel that delivers `ticks' of a clock at intervals,
// deviating with constrained random jitter.
type Ticker struct {
	C      <-chan time.Time // The channel on which the ticks are delivered.
	cancel context.CancelFunc
}

// NewTicker returns a new Ticker containing a channel that will send the
// time with a period specified by the duration argument, that will randomly
// jitter based on the provided the scaling factor.
//
// The duration d must be greater than zero; and the factor must be <= 1.0
func NewTicker(d time.Duration, factor float64) *Ticker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	if factor > 1.0 {
		panic(errors.New("factor > 1.0 for NewTicker"))
	}

	// TODO: should we allow base Context to be passed in in future?
	//
	// This could be done in a alternative constructor to keep the primary API
	// identical to time.Ticker, but would still add some noise to the public
	// methods in this package, which are currently nice and simple.
	ctx, cf := context.WithCancel(context.Background())

	// Give the channel a 1-element time buffer.
	//
	// Similar to time.Timer, if the client falls behind while reading, we will
	// need to drop ticks on the floor until the client catches up.
	c := make(chan time.Time, 1)

	go func() {
		timer := time.NewTimer(Scale(d, factor)) // initial timer
		for {
			select {
			case tc := <-timer.C:
				// Reset the internal timer for the next duration.
				//
				// Since the program has already received a value from timer.C,
				// the timer is known to have expired and the channel drained,
				// so t.Reset can be used directly.
				timer.Reset(Scale(d, factor))

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

	return &Ticker{
		C:      c,
		cancel: cf,
	}
}

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a concurrent goroutine reading from the
// channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	t.cancel()
}
