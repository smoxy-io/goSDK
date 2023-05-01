package arrays

import (
	"goSDK/util/interfaces"
	"reflect"
	"testing"
)

var unsortedStrings = []string{
	"foo",
	"baz",
	"bar",
	"apple",
	"orange",
	"banana",
}

var unsortedInts = []int{
	6,
	-1,
	0,
	7,
	8,
	2,
	4,
}

var sortedStringsAsc = []string{
	"apple",
	"banana",
	"bar",
	"baz",
	"foo",
	"orange",
}

var sortedStringsDesc = []string{
	"orange",
	"foo",
	"baz",
	"bar",
	"banana",
	"apple",
}

var sortedIntsAsc = []int{
	-1,
	0,
	2,
	4,
	6,
	7,
	8,
}

var sortedIntsDesc = []int{
	8,
	7,
	6,
	4,
	2,
	0,
	-1,
}

func TestSort(t *testing.T) {

	// default sorts ascending
	runTest(unsortedStrings, sortedStringsAsc, t)
	runTest(unsortedInts, sortedIntsAsc, t)

	// ascending option sorts ascending
	runTest(unsortedStrings, sortedStringsAsc, t, SortAscending)
	runTest(unsortedInts, sortedIntsAsc, t, SortAscending)

	// descending option sorts descending
	runTest(unsortedStrings, sortedStringsDesc, t, SortDescending)
	runTest(unsortedInts, sortedIntsDesc, t, SortDescending)

	// passing both ascending and descending sorts descending
	runTest(unsortedStrings, sortedStringsDesc, t, SortAscending, SortDescending)
	runTest(unsortedInts, sortedIntsDesc, t, SortAscending, SortDescending)
}

func runTest[A []V, V interfaces.Orderable](test A, expected A, t *testing.T, options ...SortOptions) {
	c := Clone(test)
	Sort(c, options...)

	if !reflect.DeepEqual(c, expected) {
		t.Errorf("Sort() = %v, wanted %v (options: %v)", c, expected, options)
	}
}
