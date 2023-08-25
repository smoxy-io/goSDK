package arrays

// IndexOf returns the index where value is stored in slice, or -1 if slice does not contain value
func IndexOf[A ~[]V, V comparable](slice A, value V) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}

// Contains returns true if slice contains value
func Contains[A ~[]V, V comparable](slice A, value V) bool {
	return IndexOf(slice, value) > -1
}
