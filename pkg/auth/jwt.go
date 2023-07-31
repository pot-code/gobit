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

func (j *JwtAuth) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(j.method, claims)
	return token.SignedString(j.signKey)
}

func (j *JwtAuth) Validate(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.validateKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

func (j *JwtAuth) GenerateTokenStr(data map[string]interface{}, exp time.Duration) (string, error) {
	expires := time.Now().Add(exp).Unix()
	claims := jwt.MapClaims{
		"exp": expires,
	}
	for k, v := range data {
		claims[k] = v
	}
	return j.Sign(claims)
}
