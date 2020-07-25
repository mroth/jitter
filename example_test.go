package jitter_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mroth/jitter"
)

func ExampleNewTicker() {
	// ticker with base duration of 1 second and 0.5 scaling factor
	ticker := jitter.NewTicker(time.Second, 0.5)
	defer ticker.Stop()

	prev := time.Now()
	for i := 0; i < 10; i++ {
		t := <-ticker.C // time elapsed random in range [0.5s, 1.5s]
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
