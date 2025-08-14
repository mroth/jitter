package jitter_test

import (
	"fmt"
	"time"

	"github.com/mroth/jitter"
)

func ExampleNewTicker() {
	// ticker with base duration of 10 milliseconds and 0.5 scaling factor
	ticker := jitter.NewTicker(10*time.Millisecond, 0.5)
	defer ticker.Stop()

	prev := time.Now()
	for range 5 {
		t := <-ticker.C // time elapsed is random in range (5ms, 15ms).
		fmt.Println("Time elapsed since last tick: ", t.Sub(prev))
		prev = t
	}
}
