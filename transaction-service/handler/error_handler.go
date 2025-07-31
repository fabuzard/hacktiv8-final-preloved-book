package handler

import (
	"fmt"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var status int
	var message string

	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
		message = fmt.Sprintf("%v", he.Message)
	} else {
		switch {
		case err == nil:
			return
		case err == utils.ErrUserNotFound:
			status = http.StatusNotFound
			message = err.Error()
		case err == utils.ErrUserForbidden:
			status = http.StatusForbidden
			message = err.Error()
		case err == utils.ErrBadReq:
			status = http.StatusBadRequest
			message = err.Error()
		case err == utils.ErrUnauthorized:
			status = http.StatusUnauthorized
			message = err.Error()
		default:
			status = http.StatusInternalServerError
			message = "internal server error"
		}
	}

	c.JSON(status, map[string]string{
		"error": message,
	})
}
