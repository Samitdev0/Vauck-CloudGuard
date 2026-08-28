package middleware

import (
	"net/http"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/labstack/echo/v4"
)

func RequireRole(allowedRoles ...auth.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, err := auth.GetRole(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "authenticated identity not found",
				})
			}

			if err := auth.RequireRole(role, allowedRoles...); err != nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "insufficient permissions",
				})
			}

			return next(c)
		}
	}
}
