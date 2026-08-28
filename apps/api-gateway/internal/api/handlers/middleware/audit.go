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

			auditContext := audit.FromEcho(c)

			auditLog := &models.AuditLog{
				ID:         "",
				Action:     c.Request().Method,
				Resource:   c.Path(),
				Method:     auditContext.Method,
				Path:       auditContext.Path,
				StatusCode: c.Response().Status,
				IPAddress:  nil,
				UserAgent:  auditContext.UserAgent,
				Metadata:   []byte(`{}`),
				CreatedAt:  time.Now().UTC(),
			}

			if auditContext.IP != "" {
				auditLog.IPAddress = &auditContext.IP
			}

			if auditLog.Resource == "" {
				auditLog.Resource = "unknown"
			}

			if auditContext.Identity != nil {
				auditLog.TenantID = auditContext.Identity.TenantID
				auditLog.UserID = &auditContext.Identity.UserID
			}

			metadata := map[string]string{
				"request_id": auditContext.RequestID,
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
