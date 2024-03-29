# jitter
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mroth/jitter)](https://pkg.go.dev/github.com/mroth/jitter)
[![CodeFactor](https://www.codefactor.io/repository/github/mroth/jitter/badge)](https://www.codefactor.io/repository/github/mroth/jitter)
[![Build Status](https://github.com/mroth/jitter/workflows/test/badge.svg)](https://github.com/mroth/jitter/actions)
[![codecov](https://codecov.io/gh/mroth/jitter/branch/main/graph/badge.svg)](https://codecov.io/gh/mroth/jitter)

A simple Go library providing functionality for generating durations and tickers
that deviate from true periodicity within specified bounds.

Most notably, contains a nearly API compatible version of `time.Ticker` with
definable jitter.

## Usage

For usage details, see the [Go documentation](https://pkg.go.dev/github.com/mroth/jitter).

### Example Ticker

```go
// ticker with base duration of 1 second and 0.5 scaling factor
ticker := jitter.NewTicker(time.Second, 0.5)
defer ticker.Stop()

prev := time.Now()
for {
    t := <-ticker.C // time elapsed random in range (0.5s, 1.5s)
    fmt.Println("Time elapsed since last tick: ", t.Sub(prev))
    prev = t
}
```
