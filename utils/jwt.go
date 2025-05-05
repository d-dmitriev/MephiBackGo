package utils

import (
	"bank-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var cfg = config.LoadConfig()

func GenerateJWT(userID string) (string, error) {
	duration, _ := time.ParseDuration(cfg.JWTTtl)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
