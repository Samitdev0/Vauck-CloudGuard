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

	// ==========================
	// Global Error Handler
	// ==========================

	e.HTTPErrorHandler = apierrors.HTTPErrorHandler

	// ==========================
	// Global Middlewares
	// ==========================

	e.Use(apimiddleware.RequestID)
	e.Use(apimiddleware.RequestLogger)
	e.Use(apimiddleware.Recovery())
	e.Use(apimiddleware.Security)
	e.Use(apimiddleware.CORS)
	e.Use(apimiddleware.RateLimit)

	// ==========================
	// Audit Middleware
	// ==========================

	_ = auditService

	// ==========================
	// Public Routes
	// ==========================

	routes.RegisterPublic(e)

	// ==========================
	// Health Routes
	// ==========================

	routes.RegisterHealth(e)

	// ==========================
	// API v1
	// ==========================

	apiGroup := e.Group("/api")
	v1 := apiGroup.Group("/v1")

	// ==========================
	// Authentication
	// ==========================

	routes.RegisterAuth(
		v1,
		authService,
		jwtService,
	)

	// ==========================
	// Admin
	// ==========================

	routes.RegisterAdmin(
		v1,
		jwtService,
	)

	// ==========================
	// Core API
	// ==========================

	routes.RegisterAPI(v1)

	return e
}
