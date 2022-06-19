package jitter_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mroth/jitter"
)

func ExampleNewTicker() {
	// jitter uses the global random source, seed it appropriately. for this
	// example, we seed it to a constant for predictable output.
	rand.Seed(42)

	// ticker with base duration of 10 milliseconds and 0.5 scaling factor
	ticker := jitter.NewTicker(10*time.Millisecond, 0.5)
	defer ticker.Stop()

	prev := time.Now()
	for i := 0; i < 5; i++ {
		t := <-ticker.C // time elapsed is random in range (5ms, 15ms).
		fmt.Println("Time elapsed since last tick: ", t.Sub(prev))
		prev = t
	}
}

func ExampleScale() {
	rand.Seed(1)
	for i := 0; i < 5; i++ {
		fmt.Println(jitter.Scale(time.Second, 0.5))
	}
	//Output:
	// 1.44777941s
	// 582.153551ms
	// 1.166145821s
	// 735.010051ms
	// 787.113937ms
}
