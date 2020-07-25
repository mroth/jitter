package jitter

import (
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
}

func BenchmarkScale(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scale(time.Second, 0.5)
	}
}
