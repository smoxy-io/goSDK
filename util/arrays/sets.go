package arrays

func Union[T comparable](a []T, b []T) []T {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}

	if len(a) == 0 {
		return Unique(b)
	}

	if len(b) == 0 {
		return Unique(a)
	}

	return Unique(a, b)
}

func Intersection[T comparable](a []T, b []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	t := make(map[T]bool)
	ret := make([]T, 0)

	for _, v := range a {
		if _, ok := t[v]; ok {
			// already seen this value
			continue
		}

		t[v] = true

		if !Contains(b, v) {
			continue
		}

		ret = append(ret, v)
	}

	for _, v := range b {
		if _, ok := t[v]; ok {
			// already seen this value
			continue
		}

		t[v] = true

		if !Contains(a, v) {
			continue
		}

		ret = append(ret, v)
	}

	return ret
}

func Difference[T comparable](a []T, b []T) []T {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}

	if len(a) == 0 {
		return nil
	}

	if len(b) == 0 {
		return a
	}

	ret := make([]T, 0)
	t := make(map[T]bool)

	for _, v := range a {
		if _, ok := t[v]; ok {
			// already seen this value
			continue
		}

		t[v] = true
		
		if Contains(b, v) {
			continue
		}

		ret = append(ret, v)
	}

	return ret
}

func Unique[T comparable](sets ...[]T) []T {
	if len(sets) == 0 {
		return nil
	}

	t := make(map[T]bool)
	ret := make([]T, 0)

	for _, set := range sets {
		if len(set) == 0 {
			continue
		}

		for _, v := range set {
			if _, ok := t[v]; ok {
				// already seen value
				continue
			}

			t[v] = true
			ret = append(ret, v)
		}
	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}
