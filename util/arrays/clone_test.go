package arrays

import (
	"reflect"
	"testing"
)

func TestClone(t *testing.T) {
	a := []int{
		0,
		1,
		2,
		3,
		4,
	}

	a2 := a
	a3 := Clone(a)

	if &a2[0] != &a[0] || len(a2) != len(a) {
		t.Errorf("array assignment should result in equal arrays.  %v != %v", a2, a)
	}

	if &a3[0] == &a[0] || !reflect.DeepEqual(a3, a) {
		t.Errorf("Clone() failed to create a new array or the cloned array does not contain the same values as the orignal array.  clone: %v, original: %v", a3, a)
	}
}
