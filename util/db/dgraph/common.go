package db

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	IdPrefix  = "0x"
	IdPattern = `^[0-9a-fA-F]+$`
)

var (
	IdRegexp = regexp.MustCompile(IdPattern)
)

func IdDecToHex(id string) string {
	i, iErr := strconv.ParseUint(id, 10, 64)

	if iErr != nil {
		return ""
	}

	return IdAddPrefix(strconv.FormatUint(i, 16))
}

func IdHexToDec(id string) string {
	i, iErr := strconv.ParseUint(id, 16, 64)

	if iErr != nil {
		return ""
	}

	return strconv.FormatUint(i, 10)
}

func IdTrimPrefix(id string) string {
	return strings.TrimPrefix(id, IdPrefix)
}

func IdAddPrefix(id string) string {
	return IdPrefix + id
}

func ValidId(id string, noPrefix ...bool) bool {
	if len(noPrefix) == 0 || !noPrefix[0] {
		if !strings.HasPrefix(id, IdPrefix) {
			return false
		}

		id = IdTrimPrefix(id)
	}

	return IdRegexp.MatchString(id)
}
