package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	godotenv.Load()

	jwtKey := os.Getenv("JWT_SECRET")
	fmt.Println(jwtKey)

	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid, header required")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid, bearer token is required")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid, parse token")
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("name", claims["full_name"].(string))
		c.Set("email", claims["email"].(string))
		return next(c)
	}
}
