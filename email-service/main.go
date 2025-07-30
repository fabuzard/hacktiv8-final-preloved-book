package main

import (
	"email-service/config"
	"email-service/handler"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	e.POST("/send-verification-email", handler.SendVerificationEmail)
	e.POST("/send-transaction-success", handler.SendTransactionSuccess)

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8084"))
}
