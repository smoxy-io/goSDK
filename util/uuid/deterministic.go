package uuid

import (
	"strings"

	"github.com/google/uuid"
	"github.com/smoxy-io/goSDK/util/errors"
	"github.com/zeebo/xxh3"
)

func DeterministicUUID(seed []byte) (uuid.UUID, error) {
	if len(seed) == 0 {
		return uuid.UUID{}, errors.New("empty seed")
	}

	hash := xxh3.Hash128(seed).Bytes()

	return uuid.FromBytes(hash[:])
}

func DeterministicUUIDString(seed ...string) (uuid.UUID, error) {
	return DeterministicUUID([]byte(strings.Join(seed, "")))
}
