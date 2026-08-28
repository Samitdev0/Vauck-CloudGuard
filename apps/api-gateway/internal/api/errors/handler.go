package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	problem := Problem{
		Type:   "about:blank",
		Title:  http.StatusText(code),
		Status: code,
		Detail: err.Error(),
	}

	if id, ok := c.Get("request_id").(string); ok {
		problem.RequestID = id
	}

	_ = c.JSON(code, problem)
}
