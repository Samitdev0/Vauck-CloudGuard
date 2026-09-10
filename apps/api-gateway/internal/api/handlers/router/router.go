package api

import (
	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	apierrors "github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/errors"
	apimiddleware "github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers/middleware"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers/routes"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
)

func NewRouter(
	authService *services.AuthService,
	jwtService *auth.JWT,
	auditService *services.AuditService,
) *echo.Echo {
	e := echo.New()

	// Global error handling.
	e.HTTPErrorHandler = apierrors.HTTPErrorHandler

	// Global middleware pipeline.
	e.Use(apimiddleware.RequestID)
	e.Use(apimiddleware.RequestLogger)
	e.Use(apimiddleware.Recovery())
	e.Use(apimiddleware.Security)
	e.Use(apimiddleware.CORS)
	e.Use(apimiddleware.RateLimit)
	e.Use(apimiddleware.Audit(auditService))

	// Public routes.
	routes.RegisterPublic(e)

	// Health routes.
	routes.RegisterHealth(e)

	// API v1.
	apiGroup := e.Group("/api")
	v1 := apiGroup.Group("/v1")

	// Authentication routes.
	routes.RegisterAuth(
		v1,
		authService,
		jwtService,
	)

	// Administration routes.
	routes.RegisterAdmin(
		v1,
		jwtService,
	)

	// Core API routes.
	routes.RegisterAPI(v1)

	return e
}
