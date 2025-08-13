package middleware

import (
	"book-service/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func JwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, err := helpers.ExtractToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or missing token")
			}
			c.Set("user_id", strconv.FormatUint(uint64(userID), 10))
			c.Set("role", role)
			return next(c)
		}
	}
}
