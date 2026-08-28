package middleware

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Recovery() echo.MiddlewareFunc {
	return echomiddleware.Recover()
}
