package gcm256

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/smoxy-io/goSDK/util/errors"
)

const (
	AES256KeySize = 32
)

func Encrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("no key provided")
	}

	if len(data) == 0 {
		// nothing to encrypt
		return nil, nil
	}

	block, bErr := aes.NewCipher(key)

	if bErr != nil {
		return nil, bErr
	}

	aesgcm, gErr := cipher.NewGCM(block)

	if gErr != nil {
		return nil, gErr
	}

	nonce := make([]byte, aesgcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// EncryptString returns the encrypted data as a base64 encoded string
// expects key to be base64 encoded
func EncryptString(key string, data string) (string, error) {
	k, kErr := base64.StdEncoding.DecodeString(key)

	if kErr != nil {
		return "", kErr
	}

	cipherText, cErr := Encrypt(k, []byte(data))

	if cErr != nil {
		return "", cErr
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("no key provided")
	}

	if len(data) == 0 {
		// no data to decrypt
		return nil, nil
	}

	block, bErr := aes.NewCipher(key)

	if bErr != nil {
		return nil, bErr
	}

	aesgcm, gErr := cipher.NewGCM(block)

	if gErr != nil {
		return nil, gErr
	}

	nonceSize := aesgcm.NonceSize()

	if len(data) < nonceSize {
		return nil, errors.New("invalid data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

// DecryptString expects key and encrypted data to be base64 encoded
func DecryptString(key string, data string) (string, error) {
	k, kErr := base64.StdEncoding.DecodeString(key)

	if kErr != nil {
		return "", kErr
	}

	d, dErr := base64.StdEncoding.DecodeString(data)

	if dErr != nil {
		return "", dErr
	}

	plainText, pErr := Decrypt(k, d)

	if pErr != nil {
		return "", pErr
	}

	return string(plainText), nil
}
