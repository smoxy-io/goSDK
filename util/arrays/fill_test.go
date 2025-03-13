package arrays

import (
	"slices"
	"testing"
)

type fillTest[T any] struct {
	size   int
	v      T
	expect []T
}

func TestFill(t *testing.T) {
	tests := []*fillTest[any]{
		&fillTest[any]{
			size:   0,
			v:      0,
			expect: make([]any, 0),
		},
		&fillTest[any]{
			size:   4,
			v:      0,
			expect: []any{0, 0, 0, 0},
		},
		&fillTest[any]{
			size:   3,
			v:      "foo",
			expect: []any{"foo", "foo", "foo"},
		},
	}

	for i, test := range tests {
		a := Fill(test.size, test.v)

		if !slices.Equal(test.expect, a) {
			t.Errorf("Fill[%d]: expected %v, got %v", i, test.expect, a)
		}
	}
}
