package maps

import (
	"goSDK/util/arrays"
	"goSDK/util/interfaces"
	"reflect"
	"testing"
)

var stringStringTests = []map[string]string{
	{ // test 1
		"foo": "bar",
	},
}
var intIntTests = []map[int]int{
	{ // test 2
		1: 2,
	},
}
var stringIntTests = []map[string]int{
	{ // test 3
		"baz": 1,
		"bar": 2,
	},
}
var intStringTests = []map[int]string{
	{ // test 4
		-1: "blah",
		0:  "foo",
		1:  "bar",
	},
}

func TestKeys(t *testing.T) {
	wantStrings := [][][]string{
		{ // test 1
			{
				"foo",
			},
		},
		{ // test 3
			{
				"bar",
				"baz",
			},
		},
	}
	wantInts := [][][]int{
		{ // test 2
			{
				1,
			},
		},
		{ // test 4
			{
				-1,
				0,
				1,
			},
		},
	}

	totalTests := 0

	runOrderableTests(stringStringTests, wantStrings[0], Keys[map[string]string], &totalTests, t)
	runOrderableTests(intIntTests, wantInts[0], Keys[map[int]int], &totalTests, t)
	runOrderableTests(stringIntTests, wantStrings[1], Keys[map[string]int], &totalTests, t)
	runOrderableTests(intStringTests, wantInts[1], Keys[map[int]string], &totalTests, t)
}

func TestValues(t *testing.T) {
	wantStrings := [][][]string{
		{ // test 1
			{
				"bar",
			},
		},
		{ // test 4
			{
				"bar",
				"blah",
				"foo",
			},
		},
	}
	wantInts := [][][]int{
		{ // test 2
			{
				2,
			},
		},
		{ // test 3
			{
				1,
				2,
			},
		},
	}

	totalTests := 0

	runOrderableTests(stringStringTests, wantStrings[0], Values[map[string]string], &totalTests, t)
	runOrderableTests(intIntTests, wantInts[0], Values[map[int]int], &totalTests, t)
	runOrderableTests(stringIntTests, wantInts[1], Values[map[string]int], &totalTests, t)
	runOrderableTests(intStringTests, wantStrings[1], Values[map[int]string], &totalTests, t)
}

func runOrderableTests[T ~map[K]V, K interfaces.Orderable, V interfaces.Orderable, R interfaces.Orderable](m []T, expected [][]R, fn func(T) []R, totalTests *int, t *testing.T) {
	for i, test := range m {
		res := fn(test)
		arrays.Sort(res)

		if !reflect.DeepEqual(res, expected[i]) {
			t.Errorf("Test %v: Keys() = %v, wanted %v", *totalTests+1, res, expected[i])
		}

		*totalTests++
	}
}

func TestClone(t *testing.T) {
	m := map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	}

	m2 := m
	m3 := Clone(m)

	if !reflect.DeepEqual(m2, m) {
		t.Errorf("map assignment should result in a copy by referrence")
	}

	if reflect.DeepEqual(m3, m) {
		t.Errorf("Clone() should result in a new map")
	}
}
