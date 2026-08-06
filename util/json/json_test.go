package json

import (
	"reflect"
	"strings"
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

func TestToStringPretty(t *testing.T) {
	obj := TestObj{
		Foo: "bar\nbaz",
		Bar: 4678,
		Baz: []string{"lorim", "ipsum"},
	}

	expected := strings.TrimSpace(`
{
  "foo": "bar\nbaz",
  "bar": 4678,
  "baz": [
    "lorim",
    "ipsum"
  ]
}`)

	str := ToStringPretty(obj)

	if str != expected {
		t.Errorf("ToStringPretty() =\n'%v',\nwanted:\n'%v'", str, expected)
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
