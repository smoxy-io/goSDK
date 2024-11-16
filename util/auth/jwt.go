package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	JwtContextKey = "authClaims"
)

// different types of Jwt Tokens
const (
	JwtTypePAuth = iota + 1 // primary auth token
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (m *JWTManager) Generate(user User) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"exp": now.Add(m.tokenDuration).UTC().Unix(),
		"nbf": now.Add(-1 * time.Second * 5).UTC().Unix(), // support remote clocks running a few seconds behind
		"sub": user.GetId(),
		"aud": user.GetStringRoles(),
		"iat": now.UTC().Unix(),
		"typ": JwtTypePAuth,
	})

	return token.SignedString(m.secretKey)
}

func (m *JWTManager) Verify(token string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, m.parseKeyFunc)

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*UserClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.Expired() {
		return nil, jwt.ErrTokenExpired
	}

	if claims.TimeTravel() {
		return nil, jwt.ErrTokenUsedBeforeIssued
	}

	return claims, nil
}

func (m *JWTManager) parseKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	return []byte(m.secretKey), nil
}

var jwtManager *JWTManager = NewJWTManager("foo", time.Hour)

func VerifyJwt(token string) (*UserClaims, error) {
	return jwtManager.Verify(token)
}

func GenerateJwt(user User) (string, error) {
	return jwtManager.Generate(user)
}
