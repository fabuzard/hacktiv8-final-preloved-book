package helpers

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseAndValidateEmailToken(tokenStr string) (string, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("EMAIL_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	if claims["type"] != "email_verification" {
		return "", errors.New("invalid token type")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("invalid token payload")
	}

	return email, nil
}
