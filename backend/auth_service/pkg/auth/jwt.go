package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
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

func (m *Manager) NewRefreshToken(userID int64, role string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	data := RefreshTokenData{
		UserID: userID,
		Role:   role,
		Token:  token,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(jsonData), nil
}

func (m *Manager) ParseRefreshToken(refreshToken string) (*RefreshTokenData, error) {
	jsonData, err := base64.RawURLEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, err
	}

	var data RefreshTokenData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
