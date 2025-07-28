package main

import (
	"book-service/config"
	"book-service/handler"
	"book-service/repository"
	"book-service/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	e := echo.New()

	e.Validator = config.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db := config.InitDB()

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	api := e.Group("/api/v1")

	api.Use(jwtMiddleware())

	books := api.Group("/books")
	books.POST("", bookHandler.CreateBook)
	books.GET("", bookHandler.GetAllBooks)
	books.GET("/my", bookHandler.GetMyBooks)
	books.GET("/:id", bookHandler.GetBookByID)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}

func jwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", "1")
			return next(c)
		}
	}
}

