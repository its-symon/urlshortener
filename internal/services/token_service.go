package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("SECRET_KEY") // Use env var in production

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().AddDate(1, 0, 0)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
