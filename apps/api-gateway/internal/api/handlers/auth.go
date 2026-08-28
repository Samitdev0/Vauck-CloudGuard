package handlers

import (
	"net/http"
	"strings"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

type LoginRequest struct {
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {

	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	req.TenantID = strings.TrimSpace(req.TenantID)
	req.Email = strings.TrimSpace(req.Email)

	if req.TenantID == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "tenantId, email and password are required",
		})
	}

	token, err := h.service.Login(
		c.Request().Context(),
		req.TenantID,
		req.Email,
		req.Password,
	)

	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid credentials",
			})

		case services.ErrInactiveUser:
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "user is inactive",
			})

		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"accessToken": token,
		"tokenType":   "Bearer",
	})
}

func (h *AuthHandler) Me(c echo.Context) error {

	user := c.Get("user")

	return c.JSON(http.StatusOK, user)
}
