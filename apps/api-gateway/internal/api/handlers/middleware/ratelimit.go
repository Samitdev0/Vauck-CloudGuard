package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	mu       sync.Mutex
	visitors = make(map[string]*visitor)
)

func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		ip := c.RealIP()

		mu.Lock()

		v, ok := visitors[ip]

		if !ok {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
		} else {

			if time.Since(v.lastSeen) > time.Minute {
				v.count = 0
			}

			v.count++
			v.lastSeen = time.Now()

			if v.count > 100 {
				mu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded",
				})
			}
		}

		mu.Unlock()

		return next(c)
	}
}
