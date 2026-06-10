package ed25519

import (
	"encoding/base64"

	"github.com/smoxy-io/goSDK/util/errors"
	"golang.org/x/crypto/ed25519"
)

func Sign(key ed25519.PrivateKey, data []byte) ([]byte, error) {
	if len(key) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key")
	}

	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}

	return ed25519.Sign(key, data), nil
}

// SignString signs data with key
func SignString(key ed25519.PrivateKey, data string) (string, error) {
	if len(key) != ed25519.PrivateKeySize {
		return "", errors.New("invalid private key")
	}

	if data == "" {
		return "", errors.New("no data provided")
	}

	sigBytes, sErr := Sign(key, []byte(data))

	if sErr != nil {
		return "", sErr
	}

	return base64.StdEncoding.EncodeToString(sigBytes), nil
}

func VerifySignature(key ed25519.PublicKey, data []byte, sig []byte) bool {
	if len(key) != ed25519.PublicKeySize {
		return false
	}

	if len(data) == 0 || len(sig) == 0 {
		return false
	}

	return ed25519.Verify(key, data, sig)
}

// VerifySignatureString expects sig to be base64 encoded
func VerifySignatureString(key ed25519.PublicKey, data string, sig string) bool {
	if len(key) != ed25519.PublicKeySize {
		return false
	}

	if data == "" || sig == "" {
		return false
	}

	s, sErr := base64.StdEncoding.DecodeString(sig)

	if sErr != nil {
		return false
	}

	return VerifySignature(key, []byte(data), s)
}
