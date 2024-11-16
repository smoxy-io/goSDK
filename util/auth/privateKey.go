package auth

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ed25519"
	"os"
)

var ErrUnsupportedPKCS8KeyType = errors.New("unsupported PKCS8 key type")

type PrivateKey struct {
	block *pem.Block
	key   any
}

func (k *PrivateKey) ParseKey() error {
	key, kErr := x509.ParsePKCS8PrivateKey(k.block.Bytes)

	if kErr != nil {
		return kErr
	}

	switch pk := key.(type) {
	case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
		k.key = pk
	default:
		return ErrUnsupportedPKCS8KeyType
	}

	return nil
}

func (k *PrivateKey) GetSigningMethod() jwt.SigningMethod {
	switch k.key.(type) {
	case *rsa.PrivateKey:
		return jwt.SigningMethodRS512
	case *ecdsa.PrivateKey:
		return jwt.SigningMethodES512
	case ed25519.PrivateKey:
		return jwt.SigningMethodEdDSA
	default:
		return jwt.SigningMethodNone
	}
}

func (k *PrivateKey) GetKey() any {
	return k.key
}

func (k *PrivateKey) GetPublicKey() any {
	switch pk := k.key.(type) {
	case *rsa.PrivateKey:
		return &pk.PublicKey
	case *ecdsa.PrivateKey:
		return &pk.PublicKey
	case ed25519.PrivateKey:
		return pk.Public()
	default:
		return nil
	}
}

func NewPrivateKey(pemBytes []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)

	if block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid private key pem string")
	}

	pk := &PrivateKey{
		block: block,
		key:   nil,
	}

	if err := pk.ParseKey(); err != nil {
		return nil, err
	}

	return pk, nil
}

func NewPrivateKeyFromString(pemString string) (*PrivateKey, error) {
	return NewPrivateKey([]byte(pemString))
}

func NewPrivateKeyFromFile(filePath string) (*PrivateKey, error) {
	pemBytes, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return NewPrivateKey(pemBytes)
}
