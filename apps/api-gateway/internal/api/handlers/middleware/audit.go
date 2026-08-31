package middleware

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/audit"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
)

func Audit(auditService *services.AuditService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			action := audit.ResolveAction(c)

			if !action.Auditable {
				return err
			}

			auditContext := audit.FromEcho(c)

			auditLog := &models.AuditLog{
				ID:         "",
				TenantID:   "",
				UserID:     nil,
				Action:     action.Name,
				Resource:   action.Resource,
				ResourceID: nil,
				Method:     auditContext.Method,
				Path:       auditContext.Path,
				StatusCode: c.Response().Status,
				IPAddress:  nil,
				UserAgent:  auditContext.UserAgent,
				Metadata:   []byte(`{}`),
				CreatedAt:  time.Now().UTC(),
			}

			if action.ResourceID != "" {
				resourceID := action.ResourceID
				auditLog.ResourceID = &resourceID
			}

			if auditContext.IP != "" {
				ip := auditContext.IP
				auditLog.IPAddress = &ip
			}

			metadata := map[string]any{
				"request_id": auditContext.RequestID,
				"method":     auditContext.Method,
				"path":       auditContext.Path,
				"action":     action.Name,
				"resource":   action.Resource,
				"status":     c.Response().Status,
				"authenticated": audit.IsAuthenticated(
					auditContext,
				),
			}

			if action.ResourceID != "" {
				metadata["resource_id"] = action.ResourceID
			}

			if auditContext.Identity != nil {
				tenantID := auditContext.Identity.TenantID
				userID := auditContext.Identity.UserID
				role := auditContext.Identity.Role

				auditLog.TenantID = tenantID
				auditLog.UserID = &userID

				metadata["tenant_id"] = tenantID
				metadata["user_id"] = userID
				metadata["role"] = role
			}

			if auditContext.IP != "" {
				metadata["ip_address"] = auditContext.IP
			}

			if auditContext.UserAgent != "" {
				metadata["user_agent"] = auditContext.UserAgent
			}

			if data, marshalErr := json.Marshal(metadata); marshalErr == nil {
				auditLog.Metadata = data
			}

			if auditErr := auditService.Create(
				c.Request().Context(),
				auditLog,
			); auditErr != nil {
				c.Logger().Errorf(
					"failed to create audit log: %v",
					auditErr,
				)
			}

			return err
		}
	}
}
