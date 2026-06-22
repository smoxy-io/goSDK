package gcm256

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/smoxy-io/goSDK/util/errors"
)

type Key struct {
	key []byte
}

func (k *Key) IsValid() bool {
	return len(k.key) == AES256KeySize
}

func (k *Key) String() string {
	return k.KeyString()
}

func (k *Key) Key() []byte {
	return k.key
}

func (k *Key) KeyString() string {
	return base64.StdEncoding.EncodeToString(k.key)
}

func (k *Key) Encrypt(data []byte) ([]byte, error) {
	return Encrypt(k.key, data)
}

func (k *Key) EncryptString(data string) (string, error) {
	return EncryptString(base64.StdEncoding.EncodeToString(k.key), data)
}

func (k *Key) Decrypt(data []byte) ([]byte, error) {
	return Decrypt(k.key, data)
}

func (k *Key) DecryptString(data string) (string, error) {
	return DecryptString(base64.StdEncoding.EncodeToString(k.key), data)
}

func NewKey(key []byte) *Key {
	if len(key) == AES256KeySize {
		return &Key{key: key}
	}

	// attempt to create a random key
	k := make([]byte, AES256KeySize)

	if _, err := rand.Read(k); err != nil {
		return &Key{}
	}

	return &Key{key: k}
}

func ParseKey(key string) (*Key, error) {
	k, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		return nil, err
	}

	if len(k) != AES256KeySize {
		return nil, errors.New("invalid key")
	}

	return NewKey(k), nil
}
