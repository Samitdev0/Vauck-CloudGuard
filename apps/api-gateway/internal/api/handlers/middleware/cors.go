package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CORS(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		h := c.Response().Header()

		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "*")
		h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusNoContent)
		}

		return next(c)
	}
}
