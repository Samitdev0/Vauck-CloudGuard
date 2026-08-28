package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDHeader = "X-Request-ID"

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		id := c.Request().Header.Get(RequestIDHeader)

		if id == "" {
			id = uuid.NewString()
		}

		c.Response().Header().Set(RequestIDHeader, id)

		c.Set("request_id", id)

		return next(c)
	}
}
