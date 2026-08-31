package audit

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Action struct {
	Name       string
	Resource   string
	ResourceID string
	Auditable  bool
}

func ResolveAction(c echo.Context) Action {
	method := strings.ToUpper(strings.TrimSpace(c.Request().Method))
	path := strings.TrimSpace(c.Path())

	switch {
	case path == "/health/live":
		return Action{
			Name:      "health.live",
			Resource:  "health",
			Auditable: false,
		}

	case path == "/health/ready":
		return Action{
			Name:      "health.ready",
			Resource:  "health",
			Auditable: false,
		}

	case path == "/":
		return Action{
			Name:      "system.root",
			Resource:  "system",
			Auditable: true,
		}

	case path == "/api/v1/status":
		return Action{
			Name:      "api.status",
			Resource:  "api",
			Auditable: true,
		}

	case path == "/api/v1/auth/login":
		return Action{
			Name:      "authentication.login",
			Resource:  "authentication",
			Auditable: true,
		}

	case path == "/api/v1/auth/me":
		return Action{
			Name:      "authentication.me",
			Resource:  "authentication",
			Auditable: true,
		}

	case path == "/api/v1/admin/me":
		return Action{
			Name:      "administration.me",
			Resource:  "administration",
			Auditable: true,
		}
	}

	return Action{
		Name:       methodAction(method),
		Resource:   resourceFromPath(path),
		ResourceID: resourceID(c),
		Auditable:  true,
	}
}

func methodAction(method string) string {
	switch method {
	case http.MethodGet:
		return "resource.read"

	case http.MethodPost:
		return "resource.create"

	case http.MethodPut:
		return "resource.update"

	case http.MethodPatch:
		return "resource.update"

	case http.MethodDelete:
		return "resource.delete"

	default:
		return "resource.access"
	}
}

func resourceFromPath(path string) string {
	path = strings.Trim(path, "/")

	if path == "" {
		return "system"
	}

	parts := strings.Split(path, "/")

	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "v1" {
		return parts[2]
	}

	return parts[0]
}

func resourceID(c echo.Context) string {
	for _, name := range []string{
		"id",
		"userId",
		"user_id",
		"tenantId",
		"tenant_id",
		"resourceId",
		"resource_id",
	} {
		if value := strings.TrimSpace(c.Param(name)); value != "" {
			return value
		}
	}

	return ""
}
