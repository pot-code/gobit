package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtAuth struct {
	signKey     interface{}
	validateKey interface{}
	method      jwt.SigningMethod
}

func NewJwtAuth(secret []byte, options ...JwtOption) (*JwtAuth, error) {
	auth := &JwtAuth{
		method:      jwt.SigningMethodHS256,
		signKey:     secret,
		validateKey: secret,
	}
	for _, option := range options {
		if err := option.apply(auth); err != nil {
			return nil, err
		}
	}
	return auth, nil
}

// Sign sign token
func (auth *JwtAuth) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(auth.method, claims)
	return token.SignedString(auth.signKey)
}

// Validate validate token string, error is not nil if not passed
func (auth *JwtAuth) Validate(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return auth.validateKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

// GenerateTokenStr generate token from given data
func (auth *JwtAuth) GenerateTokenStr(data map[string]interface{}, exp time.Duration) (string, error) {
	expires := time.Now().Add(exp).Unix()
	claims := jwt.MapClaims{
		"exp": expires,
	}
	for k, v := range data {
		claims[k] = v
	}
	return auth.Sign(claims)
}
