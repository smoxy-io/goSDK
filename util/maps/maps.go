package maps

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Values returns the values of the map m.
// The values will be an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// Clone returns a copy of the map m.
// The ordering of the returned map will be in an indeterminate order.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	c := make(M, len(m))

	for k, v := range m {
		c[k] = v
	}

	return c
}

// Chunk divides the map m into x parts and returns them as a list of chunks
// Useful for parallel processing of a map
// The ordering of keys between and within chunks is indeterminate
// If the chunks cannot contain the same number of keys, the last chunk will be longer
// If the map is not large enough to provide the requested number of chunks, a single chunk will be returned
func Chunk[M ~map[K]V, K comparable, V any](m M, x int) []M {
	if x == 1 {
		// one chunk is easy
		return []M{m}
	}

	l := len(m)

	if l < 2 || x >= l {
		// map not large enough to provide the requested number of chunks
		return []M{m}
	}

	chunks := make([]M, x)

	keysPerChunk := l / x

	n := 0
	i := 0
	for k, v := range m {
		if i > 0 && i%keysPerChunk == 0 && n < keysPerChunk-1 {
			// move to the next chunk
			n++
		}

		chunks[n][k] = v

		i++
	}

	return chunks
}
