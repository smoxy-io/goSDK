package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"path"
	"slices"
	"sync/atomic"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	sdkatomic "github.com/smoxy-io/goSDK/util/atomic"
	"github.com/smoxy-io/goSDK/util/env"
	"github.com/smoxy-io/goSDK/util/errors"
)

const (
	JwtContextKey = "authClaims"
)

// different types of Jwt Tokens
const (
	JwtTypePAuth = iota + 1 // primary auth token
	JwtTypeDAuth            // delegated auth token
)

type JWTGenerateOption func(claims *UserClaims)

type JWTManager struct {
	secretKey     *atomic.Pointer[PrivateKey]
	secretKeys    *sdkatomic.RcuMap[string, *PrivateKey]
	tokenDuration time.Duration
	issuer        string
	jwksEnabled   *atomic.Bool
	jwksStore     *jwkset.MemoryJWKSet
	jwksUrl       string
}

func NewJWTManager(secretKey *PrivateKey, tokenDuration time.Duration) *JWTManager {
	m := &JWTManager{
		secretKey:     new(atomic.Pointer[PrivateKey]),
		tokenDuration: tokenDuration,
		jwksEnabled:   new(atomic.Bool),
	}

	m.secretKey.Store(secretKey)

	return m
}

func (m *JWTManager) WithIssuer(issuer string) *JWTManager {
	m.issuer = issuer
	return m
}

func (m *JWTManager) WithJwksUrl(url string) *JWTManager {
	m.jwksUrl = url
	return m
}

func (m *JWTManager) EnableJwks() error {
	if m.jwksEnabled.Swap(true) {
		return nil
	}

	m.jwksStore = jwkset.NewMemoryStorage()
	m.secretKeys = sdkatomic.NewRcuMap[string, *PrivateKey]()

	key := m.secretKey.Load()

	kid, pErr := m.publishJwk(key)

	if pErr != nil {
		return pErr
	}

	m.secretKeys.Set(kid, key)

	return pErr
}

func (m *JWTManager) DisableJwks() {
	if !m.jwksEnabled.Swap(false) {
		return
	}

	m.jwksStore = nil
	m.secretKeys = nil

	return
}

func (m *JWTManager) JWKSJson(ctx context.Context) (json.RawMessage, error) {
	return m.jwksStore.JSONPublic(ctx)
}

func (m *JWTManager) publishJwk(key *PrivateKey) (string, error) {
	jwksAlg := key.GetJwkAlgorithm()

	if jwksAlg == jwkset.AlgNone {
		return "", errors.New("jwk algorithm not supported")
	}

	draft, dErr := jwkset.NewJWKFromKey(key.GetPublicKey(), jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			ALG: jwksAlg,
			USE: jwkset.UseSig,
		},
	})

	if dErr != nil {
		return "", errors.New("failed to draft jwk: %v", dErr)
	}

	kid, kErr := JWKSThumbprint(draft.Marshal())

	if kErr != nil {
		return "", errors.New("failed to thumbprint jwk: %v", kErr)
	}

	jwk, jErr := jwkset.NewJWKFromKey(key.GetPublicKey(), jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: kid,
			ALG: jwksAlg,
			USE: jwkset.UseSig,
		},
	})

	if jErr != nil {
		return "", errors.New("failed to create jwk %s: %v", kid, jErr)
	}

	if err := m.jwksStore.KeyWrite(nil, jwk); err != nil {
		return "", errors.New("failed to store jwk %s: %v", kid, err)
	}

	return kid, nil
}

func (m *JWTManager) Generate(user User) (string, error) {
	now := time.Now()

	issuer := m.issuer

	if issuer == "" {
		issuer = GetJwtIssuer()
	}

	key := m.secretKey.Load()

	claims := &UserClaims{
		Expires:   now.Add(m.tokenDuration).UTC().Unix(),
		NotBefore: now.Add(-1 * time.Second * 5).UTC().Unix(), // support remote clocks running a few seconds behind
		UserId:    user.GetId(),
		Roles:     user.GetStringRoles(),
		Issued:    now.UTC().Unix(),
		Type:      JwtTypePAuth,
		Issuer:    issuer,
		Id:        rand.Text(),
	}

	token := jwt.NewWithClaims(key.GetSigningMethod(), claims)

	return token.SignedString(key.key)
}

// GenerateDelegate generates a delegate token for an actor from the claims of an existing token. the delegate token
// cannot expand on the original token, but can restrict from the original token.
func (m *JWTManager) GenerateDelegate(user User, actor string, claims *UserClaims) (string, error) {
	if actor == "" {
		return "", errors.New("actor is required")
	}

	if claims == nil {
		return "", errors.New("claims is required")
	}

	if user.GetId() != claims.UserId {
		return "", errors.New("user id mismatch")
	}

	now := time.Now()

	issuer := m.issuer

	if issuer == "" {
		issuer = GetJwtIssuer()
	}

	key := m.secretKey.Load()

	newClaims := &UserClaims{
		Expires:   claims.Expires,
		NotBefore: now.Add(-1 * time.Second * 5).UTC().Unix(), // support remote clocks running a few seconds behind
		UserId:    user.GetId(),
		Roles:     user.GetStringRoles(),
		Issued:    now.UTC().Unix(),
		Type:      JwtTypeDAuth,
		Issuer:    issuer,
		Id:        rand.Text(),
		Act: &Actor{
			Sub: actor,
			Act: claims.Act,
		},
	}

	for _, role := range newClaims.Roles {
		if !slices.Contains(claims.Roles, role) {
			return "", errors.New("cannot assign role %s to user %s. available roles: %s", role, user.GetId(), claims.Roles)
		}
	}

	token := jwt.NewWithClaims(key.GetSigningMethod(), newClaims)

	return token.SignedString(key.key)
}

func (m *JWTManager) Verify(token string) (*UserClaims, error) {
	var jwtToken *jwt.Token
	var err error

	if m.jwksUrl != "" {
		jwtToken, err = jwt.ParseWithClaims(token, &UserClaims{}, m.jwksKeyFunc)
	} else {
		jwtToken, err = jwt.ParseWithClaims(token, &UserClaims{}, m.parseKeyFunc)
	}

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

func (m *JWTManager) parseKeyFunc(_ *jwt.Token) (any, error) {
	return m.secretKey.Load().GetPublicKey(), nil
}

func (m *JWTManager) jwksKeyFunc(token *jwt.Token) (any, error) {
	kf, kfErr := keyfunc.NewDefault([]string{m.jwksUrl})

	if kfErr != nil {
		return nil, kfErr
	}

	return kf.Keyfunc(token)
}

var jwtManager *JWTManager

func VerifyJwt(token string) (*UserClaims, error) {
	return DefaultJwtManager().Verify(token)
}

func GenerateJwt(user User) (string, error) {
	return DefaultJwtManager().Generate(user)
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

func DefaultJwtManager() *JWTManager {
	if jwtManager != nil {
		return jwtManager
	}

	jwtManager = NewJWTManager(mustGetPrivateKey(), time.Hour)

	return jwtManager
}
