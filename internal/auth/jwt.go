package auth

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

func GenerateJWT(userID string) (string, error) {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		return "", service.ErrUserMissingJWTSecret
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		return nil, nil, service.ErrUserMissingJWTSecret
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, err
	}

	return token, claims, nil
}
