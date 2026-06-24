package hash

import "testing"

func Test_Sha512(t *testing.T) {
	test := "lorim ipsum"
	expected := "6ccbb021374b1231e4d4c4cbed2c4d927f335c25d5747539e44b8dbc4402cd839e21649bff4c7faf9b02aa00df6535f2e78d7be2df80ca53dd09484a3caea222"

	result := Sha512([]byte(test))

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func Test_Sha256(t *testing.T) {
	test := "lorim ipsum"
	expected := "fc24a5c7d238d768b0b6f2067e5ded59b23a5817a582f96074637726b8a7f851"

	result := Sha256([]byte(test))

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
