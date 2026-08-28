package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHealth(e *echo.Echo) {

	e.GET("/health/live", Live)

	e.GET("/health/ready", Ready)
}

func Live(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"status":  "UP",
		"service": "api-gateway",
	})
}

func Ready(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"status":  "READY",
		"service": "api-gateway",
		"checks": map[string]any{
			"database": "pending",
			"redis":    "pending",
			"kafka":    "pending",
		},
	})
}
