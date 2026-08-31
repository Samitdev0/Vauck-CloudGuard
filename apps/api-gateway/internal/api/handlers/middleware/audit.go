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

			action := audit.ResolveAction(
				c.Request().Method,
				c.Path(),
			)

			if !action.Auditable {
				return err
			}

			auditContext := audit.FromEcho(c)

			auditLog := &models.AuditLog{
				ID:         "",
				Action:     action.Name,
				Resource:   action.Resource,
				ResourceID: nil,
				TenantID:   "",
				UserID:     nil,
				Method:     auditContext.Method,
				Path:       auditContext.Path,
				StatusCode: c.Response().Status,
				IPAddress:  nil,
				UserAgent:  auditContext.UserAgent,
				Metadata:   []byte(`{}`),
				CreatedAt:  time.Now().UTC(),
			}

			if auditContext.IP != "" {
				ip := auditContext.IP
				auditLog.IPAddress = &ip
			}

			if auditContext.Identity != nil {
				tenantID := auditContext.Identity.TenantID
				userID := auditContext.Identity.UserID

				auditLog.TenantID = tenantID
				auditLog.UserID = &userID
			}

			metadata := map[string]any{
				"request_id": auditContext.RequestID,
				"method":     auditContext.Method,
				"path":       auditContext.Path,
				"status":     c.Response().Status,
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
