package ed25519

import (
	"encoding/base64"

	"golang.org/x/crypto/ed25519"
)

type Key struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func (k *Key) String() string {
	if pk := k.PrivateKey(); pk != nil {
		return base64.StdEncoding.EncodeToString(pk)
	}

	if len(k.publicKey) == ed25519.PublicKeySize {
		return base64.StdEncoding.EncodeToString(k.publicKey)
	}

	return ""
}

func (k *Key) PrivateKey() ed25519.PrivateKey {
	if len(k.privateKey) != ed25519.PrivateKeySize {
		return nil
	}

	return k.privateKey
}

func (k *Key) PublicKey() ed25519.PublicKey {
	if len(k.publicKey) != ed25519.PublicKeySize {
		return nil
	}

	return k.publicKey
}

func (k *Key) PrivateKeyString() string {
	if pk := k.PrivateKey(); pk == nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(k.privateKey)
}

func (k *Key) PublicKeyString() string {
	if pk := k.PublicKey(); pk == nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(k.publicKey)
}

func (k *Key) Sign(data []byte) ([]byte, error) {
	return Sign(k.privateKey, data)
}

func (k *Key) SignString(data string) (string, error) {
	return SignString(k.privateKey, data)
}

func (k *Key) VerifySignature(data []byte, sig []byte) bool {
	return VerifySignature(k.publicKey, data, sig)
}

func (k *Key) VerifySignatureString(data string, sig string) bool {
	return VerifySignatureString(k.publicKey, data, sig)
}

func NewPrivateKey(pk string) *Key {
	if pk == "" {
		return nil
	}

	b, bErr := base64.StdEncoding.DecodeString(pk)

	if bErr != nil {
		return nil
	}

	if len(b) == ed25519.PrivateKeySize {
		return &Key{
			privateKey: ed25519.PrivateKey(b),
			publicKey:  ed25519.PublicKey(b[ed25519.PublicKeySize:]),
		}
	}

	if len(b) == ed25519.PublicKeySize {
		// public key is not included, generate the public key
		k := &Key{privateKey: ed25519.NewKeyFromSeed(b)}
		k.publicKey = ed25519.PublicKey(k.privateKey[ed25519.PublicKeySize:])

		return k
	}

	return nil
}

func NewPublicKey(pk string) *Key {
	if pk == "" {
		return nil
	}

	b, bErr := base64.StdEncoding.DecodeString(pk)

	if bErr != nil {
		return nil
	}

	if len(b) != ed25519.PublicKeySize {
		return nil
	}

	return &Key{publicKey: ed25519.PublicKey(b)}
}
