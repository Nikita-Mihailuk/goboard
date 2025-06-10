package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	signingKey string
}

var (
	ErrSigningKeyEmpty = errors.New("signing key is empty")
	ErrInvalidToken    = errors.New("invalid token")
)

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, ErrSigningKeyEmpty
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) ParseAccessToken(accessToken string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
