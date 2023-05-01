package arrays

// Clone returns a copy of the array a.
// The ordering of the returned array will be maintained.
func Clone[A ~[]V, V any](a A) A {
	c := make(A, len(a))

	for k, v := range a {
		c[k] = v
	}

	return c
}
