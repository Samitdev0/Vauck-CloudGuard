package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterPublic(e *echo.Echo) {

	e.GET("/", Root)
}

func Root(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"service": "CloudGuard Sandbox",
		"version": "0.1.0",
		"status":  "running",
	})
}
