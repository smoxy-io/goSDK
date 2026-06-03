package time

import (
	"math/rand"
	"time"
)

// RandDuration generates a random duration less than n in the given time unit.
//
// example:
//
//	/* random duration in seconds in the range [0, 100) */
//	RandDuration(100, time.Second)
func RandDuration(n int, unit time.Duration) time.Duration {
	return time.Duration(int64(rand.Intn(n)) * unit.Nanoseconds())
}

// RandDurationAllowNegative generates a random duration that can be positive or negative
func RandDurationAllowNegative(n int, unit time.Duration) time.Duration {
	d := RandDuration(n, unit)

	if rand.Intn(2) == 0 {
		return d
	}

	return -d
}

// RandDurationFromPoint generates a random duration using n and unit starting from the initial duration p
// this is useful for creating a randomized interval so that things are all happening at the same time
//
// WARNING: if n * unit is greater than p then this function could return a negative duration
//
// example:
//
//	/* creates a randomized 5 minute interval in the range (3, 7) minutes */
//	RandDurationFromPoint(5 * time.Minute, 120, time.Second)
func RandDurationFromPoint(p time.Duration, n int, unit time.Duration) time.Duration {
	return p + RandDurationAllowNegative(n, unit)
}
