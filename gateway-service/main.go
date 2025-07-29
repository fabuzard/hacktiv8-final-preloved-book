package main

import (
	"gateway-service/config"
	"gateway-service/handler"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Load ENV
	config.LoadEnv()

	// 2. Init Echo instance
	e := echo.New()

	// 3. Init Handler
	h := handler.NewGatewayHandler()

	// Auth endpoints
	authGroup := e.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.GET("/users/:id", h.GetUserByID)
	authGroup.PUT("/users/:id", h.UpdateUser)
	authGroup.PATCH("/users/:id", h.UpdateBalance)
	authGroup.POST("/users/verify", h.VerifyUser)
	authGroup.POST("/users/resend-verification-email", h.ResendVerificationEmail)
	// Book endpoints
	bookGroup := e.Group("/books")
	bookGroup.GET("", h.GetBooks)
	bookGroup.GET("/:id", h.GetBookByID)
	bookGroup.POST("", h.CreateBook)
	bookGroup.PUT("/:id", h.UpdateBook)
	bookGroup.DELETE("/:id", h.DeleteBook)

	// 4. Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("✅ Gateway Service running at http://localhost:" + port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
