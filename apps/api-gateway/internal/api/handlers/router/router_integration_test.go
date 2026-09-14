package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/config"
)

type routerAuditRepository struct {
	createdAudit *models.AuditLog
	createCalled bool
}

func (r *routerAuditRepository) Create(
	_ context.Context,
	audit *models.AuditLog,
) error {
	r.createCalled = true
	r.createdAudit = audit

	return nil
}

func (r *routerAuditRepository) FindByID(
	_ context.Context,
	_ string,
	_ string,
) (*models.AuditLog, error) {
	return nil, nil
}

func TestNewRouterAuditsAuthenticatedAdminRequest(t *testing.T) {
	jwtService, err := auth.NewJWT(config.JWTConfig{
		Secret: "cloudguard-test-secret-12345678901234567890",
		Issuer: "CloudGuard",
	})
	if err != nil {
		t.Fatalf(
			"failed to create JWT service: %v",
			err,
		)
	}

	repository := &routerAuditRepository{}

	auditService := services.NewAuditService(
		repository,
	)

	router := NewRouter(
		nil,
		jwtService,
		auditService,
	)

	token, err := jwtService.GenerateToken(
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440000",
		"admin@cloudguard.test",
		string(auth.RoleAdmin),
	)
	if err != nil {
		t.Fatalf(
			"failed to generate JWT: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/admin/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"User-Agent",
		"CloudGuard-Router-Test",
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
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
			"expected status %d, got %d",
			http.StatusOK,
			audit.StatusCode,
		)
	}

	if audit.UserAgent != "CloudGuard-Router-Test" {
		t.Fatalf(
			"expected user agent %q, got %q",
			"CloudGuard-Router-Test",
			audit.UserAgent,
		)
	}

	if audit.CreatedAt.IsZero() {
		t.Fatal(
			"expected audit CreatedAt to be populated",
		)
	}
}
