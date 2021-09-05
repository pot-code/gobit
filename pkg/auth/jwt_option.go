package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type JwtOption interface {
	apply(*JwtProvider) error
}

type optionFunc func(*JwtProvider) error

func (o optionFunc) apply(jp *JwtProvider) error {
	return o(jp)
}

func SetJwtMethod(method jwt.SigningMethod) JwtOption {
	return optionFunc(func(jp *JwtProvider) error {
		jp.method = method
		return nil
	})
}

func SetJwtRSAPublicKey(pem []byte) JwtOption {
	return optionFunc(func(jp *JwtProvider) error {
		pk, err := jwt.ParseRSAPublicKeyFromPEM(pem)
		if err != nil {
			return fmt.Errorf("failed to parse provided public key: %w", err)
		}
		jp.validateKey = pk
		return nil
	})
}

func SetJwtRSAPrivate(pem []byte, password string) JwtOption {
	return optionFunc(func(jp *JwtProvider) error {
		pk, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(pem, password)
		if err != nil {
			return fmt.Errorf("failed to parse provided private key: %w", err)
		}
		jp.signKey = pk
		return nil
	})
}
