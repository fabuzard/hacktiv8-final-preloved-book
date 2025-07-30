package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type GatewayHandler struct {
	AuthServiceURL string
	BookServiceURL string
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
		BookServiceURL: os.Getenv("BOOK_SERVICE_URL"),
	}
}

// Auth
func (h *GatewayHandler) Register(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/register")
}
func (h *GatewayHandler) Login(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/login")
}
func (h *GatewayHandler) GetUserByID(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}
func (h *GatewayHandler) UpdateUser(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}
func (h *GatewayHandler) UpdateBalance(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}
func (h *GatewayHandler) VerifyUser(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/verify")
}
func (h *GatewayHandler) ResendVerificationEmail(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/resend-verification-email")
}

// Books
func (h *GatewayHandler) GetBooks(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books")
}
func (h *GatewayHandler) GetBookByID(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}
func (h *GatewayHandler) CreateBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books")
}

func (h *GatewayHandler) UpdateBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}
func (h *GatewayHandler) DeleteBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}

func proxyRequest(c echo.Context, target string) error {

	targetURL := target + "?" + c.QueryParams().Encode()
	fmt.Println("🔗 Forwarding to:", targetURL)

	bodyBuf := new(bytes.Buffer)
	if _, err := io.Copy(bodyBuf, c.Request().Body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest(c.Request().Method, targetURL, bodyBuf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Forward headers (especially Authorization)
	req.Header = make(http.Header)
	for key, values := range c.Request().Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read response from service")
	}

	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
