package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	jwtIssuer = ""
)

type UserClaims struct {
	UserId    string   `json:"sub,omitempty"`
	Roles     []string `json:"aud,omitempty"`
	Expires   int64    `json:"exp,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Issued    int64    `json:"iat,omitempty"`
	//Type      int      `json:"typ,omitempty"`
	Issuer string `json:"iss,omitempty"`
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
