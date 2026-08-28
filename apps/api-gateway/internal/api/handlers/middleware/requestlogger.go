package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		start := time.Now()

		err := next(c)

		log.Info().
			Str("request_id", c.Response().Header().Get("X-Request-ID")).
			Str("method", c.Request().Method).
			Str("path", c.Path()).
			Int("status", c.Response().Status).
			Dur("latency", time.Since(start)).
			Msg("http request")

		return err
	}
}
