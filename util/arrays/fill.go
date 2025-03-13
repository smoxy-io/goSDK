package arrays

// Fill creates a new array of length size with every element in the array set to v
func Fill[T any](size int, v T) []T {
	if size < 0 {
		size = 0
	}

	ret := make([]T, size)

	for i, _ := range ret {
		ret[i] = v
	}

	return ret
}
