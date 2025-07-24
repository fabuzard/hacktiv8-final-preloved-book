package main

import (
	"auth-service/config"
	"auth-service/handler"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/routes"
	"auth-service/service"

	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()
	if db == nil {
		fmt.Println("Database connection failed, exiting...")
	}

	// Migrate the User model
	db.AutoMigrate(&models.User{})

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	routes.SetupRoutes(e, authHandler)

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8080"))
}
