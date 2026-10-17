package auth

import (
	"github.com/MicahParks/jwkset"
	"github.com/smoxy-io/goSDK/util/crypto/hash"
	"github.com/smoxy-io/goSDK/util/errors"
	"github.com/smoxy-io/goSDK/util/json"
)

type jwksThumbprint struct {
	Crv string `json:"crv,omitempty"`
	Kty string `json:"kty,omitempty"`
	E   string `json:"e,omitempty"`
	N   string `json:"n,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
}

func JWKSThumbprint(m jwkset.JWKMarshal) (string, error) {
	tp := jwksThumbprint{
		Kty: m.KTY.String(),
	}

	switch m.KTY {
	case jwkset.KtyOKP:
		tp.Crv = m.CRV.String()
		tp.X = m.X

	case jwkset.KtyEC:
		tp.Crv = m.CRV.String()
		tp.X = m.X
		tp.Y = m.Y

	case jwkset.KtyRSA:
		tp.E = m.E
		tp.N = m.N

	default:
		return "", errors.New("issuer: no thumbprint for kty %s", m.KTY.String())
	}

	return hash.Sha256(json.ToBytes(tp)), nil
}
