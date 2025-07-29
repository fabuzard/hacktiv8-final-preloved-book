package main

import (
	"main/config"
	"main/handler"
	"main/middleware"
	"main/model"
	"main/repository"
	"main/service"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	db := config.DBInit()
	godotenv.Load()
	// Run migrations
	db.AutoMigrate(&model.Transaction{})

	transRepo := repository.NewTransactionRepository(db)
	transService := service.NewTransactionService(transRepo)
	transHandler := handler.NewTransactionHandler(transService)

	e := echo.New()
	e.HTTPErrorHandler = handler.ErrorHandler
	transGroup := e.Group("/transactions")
	transGroup.Use(middleware.AuthMiddleware)
	transGroup.POST("", transHandler.CreateTransaction)
	transGroup.GET("", transHandler.GetTransaction)
	transGroup.PUT("/:trans_id", transHandler.UpdateTransactionStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))

}
