package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	apimiddleware "github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers/middleware"
)

func RegisterAdmin(
	v1 *echo.Group,
	jwtService *auth.JWT,
) {
	admin := v1.Group("/admin")

	// ==========================
	// Authentication
	// ==========================

	admin.Use(
		apimiddleware.JWT(jwtService),
	)

	// ==========================
	// Authorization
	// ==========================

	admin.Use(
		apimiddleware.RequireRole(auth.RoleAdmin),
	)

	// ==========================
	// Admin Routes
	// ==========================

	admin.GET("/me", func(c echo.Context) error {
		identity, err := auth.GetIdentity(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "authenticated identity not found",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"userId":   identity.UserID,
			"tenantId": identity.TenantID,
			"email":    identity.Email,
			"role":     identity.Role,
		})
	})
}
