package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
)

type auditIntegrationRepository struct {
	createdAudit *models.AuditLog
	createCalled bool
}

func (r *auditIntegrationRepository) Create(
	_ context.Context,
	audit *models.AuditLog,
) error {
	r.createCalled = true
	r.createdAudit = audit

	return nil
}

func (r *auditIntegrationRepository) FindByID(
	_ context.Context,
	_ string,
	_ string,
) (*models.AuditLog, error) {
	return nil, nil
}

func TestAuditMiddlewareCreatesCompleteAuditLog(t *testing.T) {
	repository := &auditIntegrationRepository{}
	auditService := services.NewAuditService(repository)

	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/admin/me",
		nil,
	)

	req.Header.Set(
		"User-Agent",
		"CloudGuard-Test-Agent",
	)

	req.Header.Set(
		"X-Forwarded-For",
		"192.168.1.100",
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetPath("/api/v1/admin/me")
	c.Set("request_id", "request-test-123")

	c.Set(auth.ContextIdentity, &auth.Identity{
		UserID:   "550e8400-e29b-41d4-a716-446655440001",
		TenantID: "550e8400-e29b-41d4-a716-446655440000",
		Email:    "admin@cloudguard.test",
		Role:     string(auth.RoleAdmin),
	})

	handler := func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]any{
				"status": "ok",
			},
		)
	}

	middleware := Audit(auditService)
	wrapped := middleware(handler)

	err := wrapped(c)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if !repository.createCalled {
		t.Fatal(
			"expected audit repository Create to be called",
		)
	}

	if repository.createdAudit == nil {
		t.Fatal(
			"expected audit log to be created",
		)
	}

	audit := repository.createdAudit

	if audit.ID == "" {
		t.Fatal(
			"expected audit ID to be generated",
		)
	}

	if audit.TenantID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf(
			"expected tenant ID %q, got %q",
			"550e8400-e29b-41d4-a716-446655440000",
			audit.TenantID,
		)
	}

	if audit.UserID == nil {
		t.Fatal(
			"expected user ID to be populated",
		)
	}

	if *audit.UserID != "550e8400-e29b-41d4-a716-446655440001" {
		t.Fatalf(
			"expected user ID %q, got %q",
			"550e8400-e29b-41d4-a716-446655440001",
			*audit.UserID,
		)
	}

	if audit.Action != "administration.me" {
		t.Fatalf(
			"expected action %q, got %q",
			"administration.me",
			audit.Action,
		)
	}

	if audit.Resource != "administration" {
		t.Fatalf(
			"expected resource %q, got %q",
			"administration",
			audit.Resource,
		)
	}

	if audit.ResourceID != nil {
		t.Fatal(
			"expected resource ID to be nil",
		)
	}

	if audit.Method != http.MethodGet {
		t.Fatalf(
			"expected method %q, got %q",
			http.MethodGet,
			audit.Method,
		)
	}

	if audit.Path != "/api/v1/admin/me" {
		t.Fatalf(
			"expected path %q, got %q",
			"/api/v1/admin/me",
			audit.Path,
		)
	}

	if audit.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status code %d, got %d",
			http.StatusOK,
			audit.StatusCode,
		)
	}

	if audit.IPAddress == nil {
		t.Fatal(
			"expected IP address to be populated",
		)
	}

	if *audit.IPAddress == "" {
		t.Fatal(
			"expected IP address to not be empty",
		)
	}

	if audit.UserAgent != "CloudGuard-Test-Agent" {
		t.Fatalf(
			"expected user agent %q, got %q",
			"CloudGuard-Test-Agent",
			audit.UserAgent,
		)
	}

	if audit.CreatedAt.IsZero() {
		t.Fatal(
			"expected CreatedAt to be generated",
		)
	}

	if audit.CreatedAt.After(time.Now().UTC()) {
		t.Fatal(
			"expected CreatedAt to not be in the future",
		)
	}

	var metadata map[string]any

	if err := json.Unmarshal(
		audit.Metadata,
		&metadata,
	); err != nil {
		t.Fatalf(
			"expected valid audit metadata JSON, got error: %v",
			err,
		)
	}

	expectedMetadata := map[string]any{
		"request_id":    "request-test-123",
		"method":        "GET",
		"path":          "/api/v1/admin/me",
		"action":        "administration.me",
		"resource":      "administration",
		"status":        float64(http.StatusOK),
		"outcome":       "success",
		"authenticated": true,
		"tenant_id":     "550e8400-e29b-41d4-a716-446655440000",
		"user_id":       "550e8400-e29b-41d4-a716-446655440001",
		"role":          string(auth.RoleAdmin),
		"user_agent":    "CloudGuard-Test-Agent",
	}

	for key, expected := range expectedMetadata {
		got, exists := metadata[key]

		if !exists {
			t.Fatalf(
				"expected metadata field %q to exist",
				key,
			)
		}

		if got != expected {
			t.Fatalf(
				"expected metadata[%q] = %v, got %v",
				key,
				expected,
				got,
			)
		}
	}

	if _, exists := metadata["duration_ms"]; !exists {
		t.Fatal(
			"expected duration_ms to exist in metadata",
		)
	}
}
