package strings

import (
	"unicode"
	"unicode/utf8"
)

// FirstToLower converts the first character in s to lower case and returns the resulting string
func FirstToLower[T ~string](s T) T {
	r, size := utf8.DecodeRuneInString(string(s))

	if r == utf8.RuneError && size <= 1 {
		return s
	}

	lc := unicode.ToLower(r)

	if r == lc {
		return s
	}

	return T(lc) + s[size:]
}

// FirstToUpper converts the first character in s to upper case and returns the resulting string
func FirstToUpper[T ~string](s T) T {
	r, size := utf8.DecodeRuneInString(string(s))

	if r == utf8.RuneError && size <= 1 {
		return s
	}

	uc := unicode.ToUpper(r)

	if r == uc {
		return s
	}

	return T(uc) + s[size:]
}
