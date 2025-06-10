package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewJWT(userID, role string, ttl time.Duration) (string, error)
	Parse(accessToken string) (*CustomClaims, error)
	NewRefreshToken() (string, error)
}

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

func NewManager(signingKey string) (TokenManager, error) {
	if signingKey == "" {
		return nil, ErrSigningKeyEmpty
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID, role string, ttl time.Duration) (string, error) {
	claims := &CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*CustomClaims, error) {
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

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
