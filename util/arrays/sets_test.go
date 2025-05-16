package arrays

import (
	"reflect"
	"testing"
)

var (
	testSetA = []string{"a", "b", "c", "d", "e"}
	testSetB = []string{"a", "b", "c", "f", "g"}
	testSetC = []string{"a", "b", "c", "d", "e", "f", "g"}
	testSetD = []string{"a", "b", "c", "d", "e", "e"}

	testSetAi = []int{1, 2, 3, 4, 5}
	testSetBi = []int{1, 2, 3, 6, 7}
	testSetCi = []int{1, 2, 3, 4, 5, 6, 7}
	testSetDi = []int{1, 2, 3, 4, 5, 5}

	testSets  = [][]string{testSetA, testSetB, testSetC}
	testSetsI = [][]int{testSetAi, testSetBi, testSetCi}
)

func TestUnique(t *testing.T) {
	if u := Unique[string](nil); u != nil {
		t.Errorf("Unique[string]() should return nil for nil input")
	}

	if u := Unique[int](nil); u != nil {
		t.Errorf("Unique[int]() should return nil for nil input")
	}

	if u := Unique([]string{}); u != nil {
		t.Errorf("Unique([]string{}) should return nil for empty input")
	}

	if u := Unique([]int{}); u != nil {
		t.Errorf("Unique([]int{}) should return nil for empty input")
	}

	if u := Unique(testSetA); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Unique() failed to return the same array for an array with unique values")
	}

	if u := Unique(testSetAi); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Unique() failed to return the same array for an array with unique values")
	}

	if u := Unique(testSetD); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Unique() failed to remove duplicate values")
	}

	if u := Unique(testSetDi); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Unique() failed to remove duplicate values")
	}

	if u := Unique(testSetA, testSetB, testSetC); !reflect.DeepEqual(u, testSetC) {
		t.Errorf("Unique() failed to remove duplicate values")
	}

	if u := Unique(testSetAi, testSetBi, testSetCi); !reflect.DeepEqual(u, testSetCi) {
		t.Errorf("Unique() failed to remove duplicate values")
	}
}

func TestUnion(t *testing.T) {
	if u := Union[string](nil, nil); u != nil {
		t.Errorf("Union[string]() should return nil for nil input")
	}

	if u := Union[int](nil, nil); u != nil {
		t.Errorf("Union[int]() should return nil for nil input")
	}

	if u := Union([]string{}, []string{}); u != nil {
		t.Errorf("Union([]string{}, []string{}) should return nil for empty input")
	}

	if u := Union([]int{}, []int{}); u != nil {
		t.Errorf("Union([]int{}, []int{}) should return nil for empty input")
	}

	if u := Union(testSetA, nil); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Union() failed to return the same array when only one set with unique values is provided")
	}

	if u := Union(nil, testSetA); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Union() failed to return the same array when only one set with unique values is provided")
	}

	if u := Union(testSetAi, nil); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Union() failed to return the same array when only one set with unique values is provided")
	}

	if u := Union(nil, testSetAi); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Union() failed to return the same array when only one set with unique values is provided")
	}

	if u := Union(testSetD, nil); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Union() failed to return an array with unique values when only one set is provided")
	}

	if u := Union(nil, testSetD); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Union() failed to return an array with unique values when only one set is provided")
	}

	if u := Union(testSetDi, nil); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Union() failed to return an array with unique values when only one set is provided")
	}

	if u := Union(nil, testSetDi); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Union() failed to return an array with unique values when only one set is provided")
	}

	if u := Union(testSetA, testSetB); !reflect.DeepEqual(u, testSetC) {
		t.Errorf("Union() failed to return union of arrays")
	}

	if u := Union(testSetAi, testSetBi); !reflect.DeepEqual(u, testSetCi) {
		t.Errorf("Union() failed to return union of arrays")
	}
}

func TestIntersection(t *testing.T) {
	if u := Intersection[string](nil, nil); u != nil {
		t.Errorf("Intersection[string]() should return nil for nil input")
	}

	if u := Intersection[int](nil, nil); u != nil {
		t.Errorf("Intersection[int]() should return nil for nil input")
	}

	if u := Intersection([]string{}, []string{}); u != nil {
		t.Errorf("Intersection([]string{}, []string{}) should return nil for empty input")
	}

	if u := Intersection([]int{}, []int{}); u != nil {
		t.Errorf("Intersection([]int{}, []int{}) should return nil for empty input")
	}

	if u := Intersection(testSetA, nil); u != nil {
		t.Errorf("Intersection() failed to return nil when only one set is provided")
	}

	if u := Intersection(nil, testSetA); u != nil {
		t.Errorf("Intersection() failed to return nil when only one set is provided")
	}

	if u := Intersection(testSetAi, nil); u != nil {
		t.Errorf("Intersection() failed to return nil when only one set is provided")
	}

	if u := Intersection(nil, testSetAi); u != nil {
		t.Errorf("Intersection() failed to return nil when only one set is provided")
	}

	if u := Intersection(testSetA, testSetB); !reflect.DeepEqual(u, []string{"a", "b", "c"}) {
		t.Errorf("Intersection() failed to return intersection of arrays")
	}

	if u := Intersection(testSetAi, testSetBi); !reflect.DeepEqual(u, []int{1, 2, 3}) {
		t.Errorf("Intersection() failed to return intersection of arrays")
	}

	if u := Intersection(testSetA, testSetD); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Intersection() failed to return intersection of arrays when one has duplicate values")
	}
}

func TestDifference(t *testing.T) {
	if u := Difference[string](nil, nil); u != nil {
		t.Errorf("Difference[string]() should return nil for nil input")
	}

	if u := Difference[int](nil, nil); u != nil {
		t.Errorf("Difference[int]() should return nil for nil input")
	}

	if u := Difference([]string{}, []string{}); u != nil {
		t.Errorf("Difference([]string{}, []string{}) should return nil for empty input")
	}

	if u := Difference([]int{}, []int{}); u != nil {
		t.Errorf("Difference([]int{}, []int{}) should return nil for empty input")
	}

	if u := Difference(testSetA, nil); !reflect.DeepEqual(u, testSetA) {
		t.Errorf("Difference() failed to return the same array when only the first set is provided")
	}

	if u := Difference(nil, testSetA); u != nil {
		t.Errorf("Difference() failed to return the same array when only the second set is provided")
	}

	if u := Difference(testSetAi, nil); !reflect.DeepEqual(u, testSetAi) {
		t.Errorf("Difference() failed to return the same array when only the first set is provided")
	}

	if u := Difference(nil, testSetAi); u != nil {
		t.Errorf("Difference() failed to return the same array when only the second set is provided")
	}

	if u := Difference(testSetA, testSetB); !reflect.DeepEqual(u, []string{"d", "e"}) {
		t.Errorf("Difference() failed to return difference of arrays")
	}

	if u := Difference(testSetAi, testSetBi); !reflect.DeepEqual(u, []int{4, 5}) {
		t.Errorf("Difference() failed to return difference of arrays")
	}

	if u := Difference(testSetD, testSetB); !reflect.DeepEqual(u, []string{"d", "e"}) {
		t.Errorf("Difference() failed to return difference of arrays when the first one has duplicate values")
	}

	if u := Difference(testSetDi, testSetBi); !reflect.DeepEqual(u, []int{4, 5}) {
		t.Errorf("Difference() failed to return difference of arrays when the first one has duplicate values")
	}
}
