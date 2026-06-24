package hash

import (
	"crypto/sha3"
	"encoding/hex"
)

func Sha512(data []byte) string {
	hash := sha3.New512()

	if _, err := hash.Write(data); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func Sha256(data []byte) string {
	hash := sha3.New256()

	if _, err := hash.Write(data); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
