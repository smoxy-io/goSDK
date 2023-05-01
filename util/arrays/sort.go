package arrays

import (
	"goSDK/util/interfaces"
	"sort"
)

type SortOptions uint16

const (
	SortAscending SortOptions = 1 << iota
	SortDescending
)

// Sort sorts the slice as specified by the passed options (default: ascending)
// The returned slice's order is deterministic
func Sort[A []V, V interfaces.Orderable](slice A, options ...SortOptions) {
	var desc bool
	var sortFn func(i, j int) bool

	for _, o := range options {
		switch o {
		case SortDescending:
			desc = true
		}
	}

	if desc {
		sortFn = func(i, j int) bool {
			return slice[i] > slice[j]
		}
	} else {
		// ascending is default
		sortFn = func(i, j int) bool {
			return slice[i] < slice[j]
		}
	}

	sort.Slice(slice, sortFn)
}
