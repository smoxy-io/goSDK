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
}

func VerifyApiKey(key ApiKey) (*UserClaims, error, int) {
	if !key.GetIsActive() {
		return nil, errors.New("using deactivated api key " + key.GetId()), http.StatusForbidden
	}

	// TODO: verify signature

	roles := make([]string, 0)

	for _, role := range key.GetRoles() {
		roles = append(roles, role.GetName())
	}

	claims := &UserClaims{
		Roles:  roles,
		Issuer: jwtIssuer,
	}

	return claims, nil, 0
}
