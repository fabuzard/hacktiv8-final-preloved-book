package helpers

import (
	"auth-service/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"role":      user.Role,
		"email":     user.Email,
		"full_name": user.Fullname,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
