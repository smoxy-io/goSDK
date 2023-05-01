package json

import (
	"reflect"
	"testing"
)

type TestObj struct {
	Foo string   `json:"foo"`
	Bar int      `json:"bar"`
	Baz []string `json:"baz"`
}

func TestToString(t *testing.T) {
	obj := TestObj{
		Foo: "bar\nbaz",
		Bar: 4678,
		Baz: []string{"lorim", "ipsum"},
	}

	expected := `{"foo":"bar\nbaz","bar":4678,"baz":["lorim","ipsum"]}`

	str := ToString(obj)

	if str != expected {
		t.Errorf("ToString() = '%v', wanted: '%v'", str, expected)
	}
}

func TestFromString(t *testing.T) {
	str := `{"foo":"bar\nbaz","bar":4678,"baz":["lorim","ipsum"]}`

	expected := TestObj{
		Foo: "bar\nbaz",
		Bar: 4678,
		Baz: []string{"lorim", "ipsum"},
	}

	obj := FromString[TestObj](str)

	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("FromString() = %v, eanted: %v", obj, expected)
	}
}
