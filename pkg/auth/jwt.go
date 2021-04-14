package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtProvider .
type JwtProvider struct {
	signKey           interface{}
	validateKey       interface{}
	Expiration        time.Duration
	RefreshExpiration time.Duration
	method            jwt.SigningMethod
}

type RSAConfig struct {
	KeyType, Password string
	Secret            []byte
}

// NewJwtProvider create a JWTUtil instance
func NewJwtProvider(
	method jwt.SigningMethod,
	secret []byte,
	expiration,
	refreshExpiration time.Duration,
	configs ...RSAConfig,
) *JwtProvider {
	var (
		signKey     interface{} = secret
		validateKey interface{} = secret
	)
	if method == jwt.SigningMethodRS256 {
		if len(configs) < 1 {
			log.Fatal("must provide a RSAConfig")
		}
		for _, config := range configs {
			if config.KeyType == "public" {
				pk, err := jwt.ParseRSAPublicKeyFromPEM(config.Secret)
				if err != nil {
					log.Fatal("failed to parse provided public key: ", err)
				}
				validateKey = pk
			} else if config.KeyType == "private" {
				pk, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(config.Secret, config.Password)
				if err != nil {
					log.Fatal("failed to parse provided private key: ", err)
				}
				signKey = pk
			} else {
				log.Fatalf("unsupported key type: '%s'", config.KeyType)
			}
		}
	}
	return &JwtProvider{
		method:            method,
		signKey:           signKey,
		validateKey:       validateKey,
		Expiration:        expiration,
		RefreshExpiration: refreshExpiration,
	}
}

// Sign sign token
func (jp *JwtProvider) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jp.method, claims)
	return token.SignedString(jp.signKey)
}

// Validate validate token string with secret and return AppTokenClaims
func (jp *JwtProvider) Validate(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jp.validateKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

// Validate validate token string with secret and return AppTokenClaims
func (jp *JwtProvider) ValidateKey() interface{} {
	return jp.validateKey
}

// GenerateTokenStr generate user token from user model
func (jp *JwtProvider) GenerateTokenStr(data map[string]interface{}) (string, error) {
	expires := time.Now().Add(jp.Expiration).Unix()
	claims := jwt.MapClaims{
		"exp": expires,
	}
	for k, v := range data {
		claims[k] = v
	}
	return jp.Sign(claims)
}

// GenerateTokenStr generate user token from user model
func (jp *JwtProvider) GenerateRefreshTokenStr(data map[string]interface{}) (string, error) {
	expires := time.Now().Add(jp.RefreshExpiration).Unix()
	claims := jwt.MapClaims{
		"exp": expires,
	}
	for k, v := range data {
		claims[k] = v
	}
	return jp.Sign(claims)
}

// GenerateTokenStr generate user token from user model
func (jp *JwtProvider) GenerateRefreshTokenCookie(name, token string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(jp.RefreshExpiration),
	}
}
