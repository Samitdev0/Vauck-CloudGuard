package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers"
	apimiddleware "github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers/middleware"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
)

func RegisterAuth(
	v1 *echo.Group,
	authService *services.AuthService,
	jwtService *auth.JWT,
) {
	authHandler := handlers.NewAuthHandler(authService)

	authGroup := v1.Group("/auth")

	// ==========================
	// Public Routes
	// ==========================

	authGroup.POST("/login", authHandler.Login)

	// ==========================
	// Protected Routes
	// ==========================

	protected := authGroup.Group("")

	protected.Use(
		apimiddleware.JWT(jwtService),
	)

	protected.GET(
		"/me",
		authHandler.Me,
	)
}
