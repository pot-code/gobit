package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtProvider .
type JwtProvider struct {
	signKey     interface{}
	validateKey interface{}
	method      jwt.SigningMethod
}

func NewJwtProvider(secret []byte, options ...JwtOption) (*JwtProvider, error) {
	jp := &JwtProvider{
		method:      jwt.SigningMethodHS256,
		signKey:     secret,
		validateKey: secret,
	}
	for _, option := range options {
		if err := option.apply(jp); err != nil {
			return nil, err
		}
	}
	return jp, nil
}

// Sign sign token
func (jp *JwtProvider) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jp.method, claims)
	return token.SignedString(jp.signKey)
}

// Validate validate token string, error is not nil if not passed
func (jp *JwtProvider) Validate(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jp.validateKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

// GenerateTokenStr generate token from given data
func (jp *JwtProvider) GenerateTokenStr(data map[string]interface{}, exp time.Duration) (string, error) {
	expires := time.Now().Add(exp).Unix()
	claims := jwt.MapClaims{
		"exp": expires,
	}
	for k, v := range data {
		claims[k] = v
	}
	return jp.Sign(claims)
}
