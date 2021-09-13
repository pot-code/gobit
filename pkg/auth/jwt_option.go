package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type JwtOption interface {
	apply(*JwtAuth) error
}

type optionFunc func(*JwtAuth) error

func (o optionFunc) apply(jp *JwtAuth) error {
	return o(jp)
}

func WithJwtMethod(method jwt.SigningMethod) JwtOption {
	return optionFunc(func(jp *JwtAuth) error {
		jp.method = method
		return nil
	})
}

func WithJwtRSAPublicKey(pem []byte) JwtOption {
	return optionFunc(func(jp *JwtAuth) error {
		pk, err := jwt.ParseRSAPublicKeyFromPEM(pem)
		if err != nil {
			return fmt.Errorf("failed to parse provided public key: %w", err)
		}
		jp.validateKey = pk
		return nil
	})
}

func WithJwtRSAPrivateKey(pem []byte, password string) JwtOption {
	return optionFunc(func(jp *JwtAuth) error {
		pk, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(pem, password)
		if err != nil {
			return fmt.Errorf("failed to parse provided private key: %w", err)
		}
		jp.signKey = pk
		return nil
	})
}
