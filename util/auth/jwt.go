package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/smoxy-io/goSDK/util/env"

	"path"
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
	secretKey     *PrivateKey
	tokenDuration time.Duration
}

func NewJWTManager(secretKey *PrivateKey, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (m *JWTManager) Generate(user User) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(m.secretKey.GetSigningMethod(), &UserClaims{
		Expires:   now.Add(m.tokenDuration).UTC().Unix(),
		NotBefore: now.Add(-1 * time.Second * 5).UTC().Unix(), // support remote clocks running a few seconds behind
		UserId:    user.GetId(),
		Roles:     user.GetStringRoles(),
		Issued:    now.UTC().Unix(),
		Type:      JwtTypePAuth,
		Issuer:    GetJwtIssuer(),
	})

	return token.SignedString(m.secretKey.key)
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
	return m.secretKey.GetPublicKey(), nil
}

var jwtManager *JWTManager = NewJWTManager(mustGetPrivateKey(), time.Hour)

func VerifyJwt(token string) (*UserClaims, error) {
	return jwtManager.Verify(token)
}

func GenerateJwt(user User) (string, error) {
	return jwtManager.Generate(user)
}

// mustGetPrivateKey returns the PrivateKey to use for JWT tokens. panics if an error occurs
func mustGetPrivateKey() *PrivateKey {
	filePath := getPrivateKeyPath()

	pk, err := NewPrivateKeyFromFile(filePath)

	if err != nil {
		panic("failed to load private key at: " + filePath + "\nerror: " + err.Error())
	}

	return pk
}

func getPrivateKeyPath() string {
	filePath := env.Get(JwtPrivateKey, DefaultJwtPrivateKey)

	if filePath == "" {
		return ""
	}

	if path.IsAbs(filePath) {
		return filePath
	}

	return path.Join(env.GetPwd(), filePath)
}
