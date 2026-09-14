package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/config"
)

type routerAuditRepository struct {
	createdAudits []*models.AuditLog
}

func (r *routerAuditRepository) Create(
	_ context.Context,
	audit *models.AuditLog,
) error {
	r.createdAudits = append(
		r.createdAudits,
		audit,
	)

	return nil
}

func (r *routerAuditRepository) FindByID(
	_ context.Context,
	_ string,
	_ string,
) (*models.AuditLog, error) {
	return nil, nil
}

func newRouterTestJWT(t *testing.T) *auth.JWT {
	t.Helper()

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

	return jwtService
}

func newRouterTest(
	t *testing.T,
) (*auth.JWT, *routerAuditRepository, *echo.Echo) {
	t.Helper()

	jwtService := newRouterTestJWT(t)

	repository := &routerAuditRepository{}

	auditService := services.NewAuditService(
		repository,
	)

	router := NewRouter(
		nil,
		jwtService,
		auditService,
	)

	return jwtService, repository, router
}

func generateRouterTestToken(
	t *testing.T,
	jwtService *auth.JWT,
	userID string,
	tenantID string,
	email string,
	role auth.Role,
) string {
	t.Helper()

	token, err := jwtService.GenerateToken(
		userID,
		tenantID,
		email,
		string(role),
	)
	if err != nil {
		t.Fatalf(
			"failed to generate JWT: %v",
			err,
		)
	}

	return token
}

func TestNewRouterAuditsAuthenticatedAdminRequest(t *testing.T) {
	jwtService, repository, router := newRouterTest(t)

	token := generateRouterTestToken(
		t,
		jwtService,
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440000",
		"admin@cloudguard.test",
		auth.RoleAdmin,
	)

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

	if len(repository.createdAudits) != 1 {
		t.Fatalf(
			"expected 1 audit log, got %d",
			len(repository.createdAudits),
		)
	}

	audit := repository.createdAudits[0]

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
			"expected audit status %d, got %d",
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

	var metadata map[string]any

	if err := json.Unmarshal(
		audit.Metadata,
		&metadata,
	); err != nil {
		t.Fatalf(
			"failed to decode audit metadata: %v",
			err,
		)
	}

	if metadata["outcome"] != "success" {
		t.Fatalf(
			"expected audit outcome %q, got %v",
			"success",
			metadata["outcome"],
		)
	}

	if metadata["authenticated"] != true {
		t.Fatalf(
			"expected authenticated metadata to be true, got %v",
			metadata["authenticated"],
		)
	}
}

func TestNewRouterRejectsUnauthenticatedAdminRequest(t *testing.T) {
	_, repository, router := newRouterTest(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/admin/me",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}

	if len(repository.createdAudits) != 0 {
		t.Fatalf(
			"expected no audit log for unauthenticated request, got %d",
			len(repository.createdAudits),
		)
	}
}

func TestNewRouterRejectsNonAdminRequest(t *testing.T) {
	jwtService, repository, router := newRouterTest(t)

	token := generateRouterTestToken(
		t,
		jwtService,
		"550e8400-e29b-41d4-a716-446655440002",
		"550e8400-e29b-41d4-a716-446655440000",
		"user@cloudguard.test",
		auth.RoleUser,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/admin/me",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}

	if len(repository.createdAudits) != 1 {
		t.Fatalf(
			"expected 1 audit log, got %d",
			len(repository.createdAudits),
		)
	}

	audit := repository.createdAudits[0]

	if audit.TenantID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf(
			"expected tenant ID %q, got %q",
			"550e8400-e29b-41d4-a716-446655440000",
			audit.TenantID,
		)
	}

	if audit.StatusCode != http.StatusForbidden {
		t.Fatalf(
			"expected audit status %d, got %d",
			http.StatusForbidden,
			audit.StatusCode,
		)
	}

	var metadata map[string]any

	if err := json.Unmarshal(
		audit.Metadata,
		&metadata,
	); err != nil {
		t.Fatalf(
			"failed to decode audit metadata: %v",
			err,
		)
	}

	if metadata["outcome"] != "denied" {
		t.Fatalf(
			"expected audit outcome %q, got %v",
			"denied",
			metadata["outcome"],
		)
	}

	if metadata["authenticated"] != true {
		t.Fatalf(
			"expected authenticated metadata to be true, got %v",
			metadata["authenticated"],
		)
	}
}
