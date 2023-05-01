package events

import "testing"

func TestRoutingKey_IsValid(t *testing.T) {
	validKeys := []RoutingKey{
		"a",
		"_",
		"1",
		"a.b.c",
		"a1_3-2.foo.blah.12-",
	}

	invalidKeys := []RoutingKey{
		".",
		"-",
		"-.#",
		"*",
		"abc..def",
		"123.-",
	}

	runValidRoutingKeyTest(validKeys, true, t)
	runValidRoutingKeyTest(invalidKeys, false, t)
}

func runValidRoutingKeyTest(tests []RoutingKey, expected bool, t *testing.T) {
	for _, test := range tests {
		if v := test.IsValid(); v != expected {
			t.Errorf("'%v'.IsValid() = %v, wanted %v", test, v, expected)
		}
	}
}

func TestRoutingKey_String(t *testing.T) {
	tests := map[RoutingKey]string{
		RoutingKey("foo"):                     "foo",
		RoutingKey("foo.bar"):                 "foo.bar",
		RoutingKey("foo.bar_baz"):             "foo.bar_baz",
		RoutingKey("foo.bar.baz"):             "foo.bar.baz",
		RoutingKey("foo.bar-baz.lorim_ipsum"): "foo.bar-baz.lorim_ipsum",
	}

	for rk, str := range tests {
		if rstr := rk.String(); rstr != str {
			t.Errorf("RoutingKey.String() = '%v', wanted: '%v'", rstr, str)
		}
	}
}
