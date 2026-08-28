package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterAPI(g *echo.Group) {

	g.GET("/status", Status)
}

func Status(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"api":     "v1",
		"message": "CloudGuard API Online",
	})
}
