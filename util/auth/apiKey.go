package auth

import (
	"errors"
	"net/http"
)

type ApiKeyRole interface {
	GetName() string
}

type ApiKey interface {
	GetId() string
	GetIsActive() bool
	GetRoles() []ApiKeyRole
	GetUser() User
}

func VerifyApiKey(key ApiKey) (*UserClaims, error, int) {
	if !key.GetIsActive() {
		return nil, errors.New("using deactivated api key " + key.GetId()), http.StatusForbidden
	}

	// TODO: verify signature

	r := make([]string, 0)

	for _, role := range key.GetRoles() {
		r = append(r, role.GetName())
	}

	claims := &UserClaims{
		Roles:  r,
		Issuer: jwtIssuer,
	}

	return claims, nil, 0
}
