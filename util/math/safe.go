package math

import "math"

// SafeUint64ToInt64 converts a uint64 to an int64, preventing overflow
func SafeUint64ToInt64(u uint64) int64 {
	if u > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(u)
}
