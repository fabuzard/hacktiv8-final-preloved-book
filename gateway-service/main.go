package main

import (
	"gateway-service/config"
	_ "gateway-service/docs"
	"gateway-service/handler"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Used Book Marketplace API Gateway
// @version 1.0
// @description This is the API Gateway for the Used Book Marketplace microservices
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. Load ENV
	config.LoadEnv()

	// 2. Init Echo instance
	e := echo.New()
	
	// Add CORS middleware
	e.Use(middleware.CORS())
	
	// Add Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

	// Transaction endpoints
	transactionGroup := e.Group("/transactions")
	transactionGroup.POST("", h.CreateTransaction)
	transactionGroup.GET("", h.GetTransactions)
	transactionGroup.PUT("/:trans_id", h.UpdateTransactionStatus)

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
