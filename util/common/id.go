package common

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"runtime/debug"
	"strconv"
)

const (
	IdSeparator  = ":"
	IdNumberBase = 36

	MaxUintBytes  = 8
	SumPrefixLen  = 4
	SumPostfixLen = 4
)

func GenerateId(prefix string, uniqueId string) (string, error) {
	h := sha256.New()

	if _, err := io.WriteString(h, prefix); err != nil {
		return "", err
	}

	if _, err := io.WriteString(h, uniqueId); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	sum = append(sum[:SumPrefixLen], sum[len(sum)-SumPostfixLen:]...)

	if len(sum) != MaxUintBytes {
		// this should never happen except by coding error
		panic([]any{"invalid id hash sum", debug.Stack()})
	}

	id, err := strconv.ParseUint(hex.EncodeToString(sum), 16, 64)

	if err != nil {
		return "", err
	}

	if prefix != "" {
		prefix = prefix + IdSeparator
	}

	return prefix + strconv.FormatUint(id, IdNumberBase), nil
}
