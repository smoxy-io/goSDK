package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtIssuer = ""
)

type Actor struct {
	Sub string `json:"sub,omitempty"`
	Act *Actor `json:"act,omitempty"`
}

type UserClaims struct {
	UserId    string            `json:"sub,omitempty"`
	Roles     []string          `json:"aud,omitempty"`
	Expires   int64             `json:"exp,omitempty"`
	NotBefore int64             `json:"nbf,omitempty"`
	Issued    int64             `json:"iat,omitempty"`
	Type      int               `json:"typ,omitempty"`
	Issuer    string            `json:"iss,omitempty"`
	Act       *Actor            `json:"act,omitempty"`
	Id        string            `json:"jti,omitempty"`
	Scope     string            `json:"scope,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

func (u *UserClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: u.ExpireTime(),
	}, nil
}

func (u *UserClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: u.IssuedTime(),
	}, nil
}

func (u *UserClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: u.NotBeforeTime(),
	}, nil
}

func (u *UserClaims) GetIssuer() (string, error) {
	return u.Issuer, nil
}

func (u *UserClaims) GetSubject() (string, error) {
	return u.UserId, nil
}

func (u *UserClaims) GetAudience() (jwt.ClaimStrings, error) {
	return u.Roles, nil
}

func (u *UserClaims) Expired() bool {
	return time.Now().After(u.ExpireTime())
}

func (u *UserClaims) TimeTravel() bool {
	return time.Now().Before(u.NotBeforeTime())
}

func (u *UserClaims) ContextWithClaims(ctx context.Context) context.Context {
	return context.WithValue(ctx, JwtContextKey, u)
}

func (u *UserClaims) ExpireTime() time.Time {
	return time.Unix(u.Expires, 0)
}

func (u *UserClaims) IssuedTime() time.Time {
	return time.Unix(u.Issued, 0)
}

func (u *UserClaims) NotBeforeTime() time.Time {
	return time.Unix(u.NotBefore, 0)
}

func (u *UserClaims) GetActor() string {
	if u.Act == nil {
		return ""
	}

	return u.Act.Sub
}

func (u *UserClaims) GetActors() []string {
	if u.Act == nil {
		return nil
	}

	actors := []string{
		u.Act.Sub,
	}

	if u.Act.Act == nil {
		return actors
	}

	a := u.Act.Act

	for a != nil {
		actors = append(actors, a.Sub)
		a = a.Act
	}

	return actors
}

func (u *UserClaims) GetScope() string {
	return u.Scope
}

func (u *UserClaims) GetExtra() map[string]string {
	return u.Extra
}

func GetUserIdFromCtx(ctx context.Context) string {
	if claims := GetUserClaimsFromCtx(ctx); claims != nil {
		return claims.UserId
	}

	return ""
}

func GetUserClaimsFromCtx(ctx context.Context) *UserClaims {
	if claims, ok := ctx.Value(JwtContextKey).(*UserClaims); ok {
		return claims
	}

	return nil
}

func SetJwtIssuer(issuer string) {
	jwtIssuer = issuer
}

func GetJwtIssuer() string {
	return jwtIssuer
}
