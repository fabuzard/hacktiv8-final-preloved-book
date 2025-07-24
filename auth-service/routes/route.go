package routes

import (
	"auth-service/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, h *handler.AuthHandler) {
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/users/:id", h.GetUserByID)
	e.PUT("/users/:id", h.UpdateUser)
}
