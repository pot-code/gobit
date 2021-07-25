package auth

import (
	"fmt"
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

type RSAKeyType int

const (
	RSAPrivate RSAKeyType = iota
	RSApublic
)

type RSAConfig struct {
	KeyType  RSAKeyType
	Password string
	Secret   []byte
}

// NewJwtProvider create a JWTUtil instance
func NewJwtProvider(
	method jwt.SigningMethod,
	secret []byte,
	expiration,
	refreshExpiration time.Duration,
	configs ...RSAConfig,
) (*JwtProvider, error) {
	var (
		signKey     interface{} = secret
		validateKey interface{} = secret
	)
	if method == jwt.SigningMethodRS256 {
		return createRSAProvider(method, secret, expiration, refreshExpiration, configs...)
	}
	return &JwtProvider{
		method:            method,
		signKey:           signKey,
		validateKey:       validateKey,
		Expiration:        expiration,
		RefreshExpiration: refreshExpiration,
	}, nil
}

func createRSAProvider(
	method jwt.SigningMethod,
	secret []byte,
	expiration,
	refreshExpiration time.Duration,
	configs ...RSAConfig,
) (*JwtProvider, error) {
	var (
		signKey     interface{} = secret
		validateKey interface{} = secret
	)
	for _, config := range configs {
		if config.KeyType == RSApublic {
			pk, err := jwt.ParseRSAPublicKeyFromPEM(config.Secret)
			if err != nil {
				return nil, fmt.Errorf("failed to parse provided public key: %w", err)
			}
			validateKey = pk
		} else if config.KeyType == RSAPrivate {
			pk, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(config.Secret, config.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to parse provided private key: %w", err)
			}
			signKey = pk
		}
	}
	return &JwtProvider{
		method:            method,
		signKey:           signKey,
		validateKey:       validateKey,
		Expiration:        expiration,
		RefreshExpiration: refreshExpiration,
	}, nil
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
