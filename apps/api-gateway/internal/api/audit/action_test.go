package audit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestResolveAction(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		path          string
		params        map[string]string
		expectedName  string
		expectedRes   string
		expectedID    string
		expectedAudit bool
	}{
		{
			name:          "health live is not auditable",
			method:        http.MethodGet,
			path:          "/health/live",
			expectedName:  "health.live",
			expectedRes:   "health",
			expectedAudit: false,
		},
		{
			name:          "health ready is not auditable",
			method:        http.MethodGet,
			path:          "/health/ready",
			expectedName:  "health.ready",
			expectedRes:   "health",
			expectedAudit: false,
		},
		{
			name:          "root",
			method:        http.MethodGet,
			path:          "/",
			expectedName:  "system.root",
			expectedRes:   "system",
			expectedAudit: true,
		},
		{
			name:          "api status",
			method:        http.MethodGet,
			path:          "/api/v1/status",
			expectedName:  "api.status",
			expectedRes:   "api",
			expectedAudit: true,
		},
		{
			name:          "authentication login",
			method:        http.MethodPost,
			path:          "/api/v1/auth/login",
			expectedName:  "authentication.login",
			expectedRes:   "authentication",
			expectedAudit: true,
		},
		{
			name:          "authentication me",
			method:        http.MethodGet,
			path:          "/api/v1/auth/me",
			expectedName:  "authentication.me",
			expectedRes:   "authentication",
			expectedAudit: true,
		},
		{
			name:          "administration me",
			method:        http.MethodGet,
			path:          "/api/v1/admin/me",
			expectedName:  "administration.me",
			expectedRes:   "administration",
			expectedAudit: true,
		},
		{
			name:          "fallback create resource",
			method:        http.MethodPost,
			path:          "/api/v1/users",
			expectedName:  "resource.create",
			expectedRes:   "users",
			expectedAudit: true,
		},
		{
			name:          "fallback update resource",
			method:        http.MethodPut,
			path:          "/api/v1/users/123",
			params:        map[string]string{"id": "123"},
			expectedName:  "resource.update",
			expectedRes:   "users",
			expectedID:    "123",
			expectedAudit: true,
		},
		{
			name:          "fallback patch resource",
			method:        http.MethodPatch,
			path:          "/api/v1/users/123",
			params:        map[string]string{"id": "123"},
			expectedName:  "resource.update",
			expectedRes:   "users",
			expectedID:    "123",
			expectedAudit: true,
		},
		{
			name:          "fallback delete resource",
			method:        http.MethodDelete,
			path:          "/api/v1/users/123",
			params:        map[string]string{"id": "123"},
			expectedName:  "resource.delete",
			expectedRes:   "users",
			expectedID:    "123",
			expectedAudit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			if len(tt.params) > 0 {
				names := make([]string, 0, len(tt.params))
				values := make([]string, 0, len(tt.params))

				for name, value := range tt.params {
					names = append(names, name)
					values = append(values, value)
				}

				c.SetParamNames(names...)
				c.SetParamValues(values...)
			}

			result := ResolveAction(c)

			if result.Name != tt.expectedName {
				t.Fatalf(
					"expected action %q, got %q",
					tt.expectedName,
					result.Name,
				)
			}

			if result.Resource != tt.expectedRes {
				t.Fatalf(
					"expected resource %q, got %q",
					tt.expectedRes,
					result.Resource,
				)
			}

			if result.ResourceID != tt.expectedID {
				t.Fatalf(
					"expected resource id %q, got %q",
					tt.expectedID,
					result.ResourceID,
				)
			}

			if result.Auditable != tt.expectedAudit {
				t.Fatalf(
					"expected auditable %v, got %v",
					tt.expectedAudit,
					result.Auditable,
				)
			}
		})
	}
}

func TestMethodAction(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{
			method: http.MethodGet,
			want:   "resource.read",
		},
		{
			method: http.MethodPost,
			want:   "resource.create",
		},
		{
			method: http.MethodPut,
			want:   "resource.update",
		},
		{
			method: http.MethodPatch,
			want:   "resource.update",
		},
		{
			method: http.MethodDelete,
			want:   "resource.delete",
		},
		{
			method: http.MethodHead,
			want:   "resource.access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			got := methodAction(tt.method)

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestResourceFromPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			path: "/",
			want: "system",
		},
		{
			path: "/api/v1/users",
			want: "users",
		},
		{
			path: "/api/v1/users/123",
			want: "users",
		},
		{
			path: "/metrics",
			want: "metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := resourceFromPath(tt.path)

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestResourceID(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name   string
		params map[string]string
		want   string
	}{
		{
			name: "id",
			params: map[string]string{
				"id": "123",
			},
			want: "123",
		},
		{
			name: "user id",
			params: map[string]string{
				"userId": "user-123",
			},
			want: "user-123",
		},
		{
			name: "tenant id",
			params: map[string]string{
				"tenantId": "tenant-123",
			},
			want: "tenant-123",
		},
		{
			name: "resource id",
			params: map[string]string{
				"resourceId": "resource-123",
			},
			want: "resource-123",
		},
		{
			name:   "no id",
			params: map[string]string{},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/test",
				nil,
			)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if len(tt.params) > 0 {
				names := make([]string, 0, len(tt.params))
				values := make([]string, 0, len(tt.params))

				for name, value := range tt.params {
					names = append(names, name)
					values = append(values, value)
				}

				c.SetParamNames(names...)
				c.SetParamValues(values...)
			}

			got := resourceID(c)

			if got != tt.want {
				t.Fatalf(
					"expected %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}
