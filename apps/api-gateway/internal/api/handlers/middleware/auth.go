package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
)

func JWT(jwtService *auth.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := strings.TrimSpace(
				c.Request().Header.Get("Authorization"),
			)

			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing token",
				})
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(header, prefix) {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization scheme",
				})
			}

			tokenString := strings.TrimSpace(
				strings.TrimPrefix(header, prefix),
			)

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing token",
				})
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
			}

			identity, err := auth.IdentityFromClaims(claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "authenticated identity not found",
				})
			}

			c.Set(auth.ContextIdentity, identity)

			c.Set("user", claims)
			c.Set("userId", identity.UserID)
			c.Set("tenantId", identity.TenantID)
			c.Set("email", identity.Email)
			c.Set("role", identity.Role)

			return next(c)
		}
	}
}
