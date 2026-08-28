package middleware

import "github.com/labstack/echo/v4"

func Security(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		h := c.Response().Header()

		h.Set("X-Frame-Options", "DENY")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("X-XSS-Protection", "1; mode=block")

		return next(c)
	}
}
