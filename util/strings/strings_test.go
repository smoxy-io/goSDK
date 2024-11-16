package strings

import "testing"

func TestFirstToLower(t *testing.T) {
	tests := []struct {
		Test string
		Want string
	}{
		{"foo", "foo"},
		{"Bar", "bar"},
		{"BAZ", "bAZ"},
	}

	for _, test := range tests {
		if got := FirstToLower(test.Test); got != test.Want {
			t.Errorf("FirstToLower(%q) = %q, want %q", test.Test, got, test.Want)
		}
	}
}

func TestFirstToUpper(t *testing.T) {
	tests := []struct {
		Test string
		Want string
	}{
		{"foo", "Foo"},
		{"Bar", "Bar"},
		{"BAZ", "BAZ"},
	}

	for _, test := range tests {
		if got := FirstToUpper(test.Test); got != test.Want {
			t.Errorf("FirstToUpper(%q) = %q, want %q", test.Test, got, test.Want)
		}
	}
}
