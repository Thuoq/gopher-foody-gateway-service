package jwt

import (
	"gopher-foody-gateway-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
}

type manager struct {
	config *config.Config
}

func NewManager(config *config.Config) TokenManager {
	return &manager{
		config: config,
	}
}

type AccessTokenClaims struct {
	PublicUserId string `json:"public_user_id"`
	SessionID    string `json:"session_id"`
	jwt.RegisteredClaims
}

func (m *manager) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.JWT.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
