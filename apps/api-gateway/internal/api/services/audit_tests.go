package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
)

type fakeAuditRepository struct {
	createCalled bool
	createdAudit *models.AuditLog

	createErr error

	findCalled   bool
	findTenantID string
	findAuditID  string
	findResult   *models.AuditLog
	findErr      error
}

func (f *fakeAuditRepository) Create(
	_ context.Context,
	audit *models.AuditLog,
) error {
	f.createCalled = true
	f.createdAudit = audit

	return f.createErr
}

func (f *fakeAuditRepository) FindByID(
	_ context.Context,
	tenantID string,
	id string,
) (*models.AuditLog, error) {
	f.findCalled = true
	f.findTenantID = tenantID
	f.findAuditID = id

	return f.findResult, f.findErr
}

func TestAuditServiceCreate(t *testing.T) {
	tests := []struct {
		name        string
		audit       *models.AuditLog
		repository  *fakeAuditRepository
		expectedErr error
	}{
		{
			name:        "nil audit log",
			audit:       nil,
			repository:  &fakeAuditRepository{},
			expectedErr: ErrAuditLogNil,
		},
		{
			name: "missing tenant id",
			audit: &models.AuditLog{
				Action:   "authentication.login",
				Resource: "authentication",
			},
			repository:  &fakeAuditRepository{},
			expectedErr: ErrTenantIDRequired,
		},
		{
			name: "missing action",
			audit: &models.AuditLog{
				TenantID: "tenant-123",
				Resource: "authentication",
			},
			repository:  &fakeAuditRepository{},
			expectedErr: ErrActionRequired,
		},
		{
			name: "missing resource",
			audit: &models.AuditLog{
				TenantID: "tenant-123",
				Action:   "authentication.login",
			},
			repository:  &fakeAuditRepository{},
			expectedErr: ErrResourceRequired,
		},
		{
			name: "repository error",
			audit: &models.AuditLog{
				TenantID: "tenant-123",
				Action:   "authentication.login",
				Resource: "authentication",
			},
			repository: &fakeAuditRepository{
				createErr: errors.New("database unavailable"),
			},
			expectedErr: errors.New("create audit log: database unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuditService(tt.repository)

			err := service.Create(
				context.Background(),
				tt.audit,
			)

			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.expectedErr)
			}

			if err.Error() != tt.expectedErr.Error() {
				t.Fatalf(
					"expected error %q, got %q",
					tt.expectedErr.Error(),
					err.Error(),
				)
			}
		})
	}
}

func TestAuditServiceCreateGeneratesIDAndTimestamp(t *testing.T) {
	repository := &fakeAuditRepository{}

	service := NewAuditService(repository)

	audit := &models.AuditLog{
		TenantID: "tenant-123",
		Action:   "authentication.login",
		Resource: "authentication",
	}

	before := time.Now().UTC()

	err := service.Create(
		context.Background(),
		audit,
	)

	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repository.createCalled {
		t.Fatal("expected repository Create to be called")
	}

	if audit.ID == "" {
		t.Fatal("expected audit ID to be generated")
	}

	if audit.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be generated")
	}

	if audit.CreatedAt.Before(before) || audit.CreatedAt.After(after) {
		t.Fatalf(
			"expected CreatedAt between %v and %v, got %v",
			before,
			after,
			audit.CreatedAt,
		)
	}

	if repository.createdAudit != audit {
		t.Fatal("expected repository to receive the same audit log")
	}
}

func TestAuditServiceCreatePreservesExistingIDAndTimestamp(t *testing.T) {
	repository := &fakeAuditRepository{}

	service := NewAuditService(repository)

	createdAt := time.Date(
		2026,
		time.January,
		1,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	audit := &models.AuditLog{
		ID:        "existing-id",
		TenantID:  "tenant-123",
		Action:    "authentication.login",
		Resource:  "authentication",
		CreatedAt: createdAt,
	}

	err := service.Create(
		context.Background(),
		audit,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if audit.ID != "existing-id" {
		t.Fatalf(
			"expected existing ID to be preserved, got %q",
			audit.ID,
		)
	}

	if !audit.CreatedAt.Equal(createdAt) {
		t.Fatalf(
			"expected existing CreatedAt to be preserved, got %v",
			audit.CreatedAt,
		)
	}
}

func TestAuditServiceFindByID(t *testing.T) {
	tenantID := "550e8400-e29b-41d4-a716-446655440000"
	auditID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	expectedAudit := &models.AuditLog{
		ID:        auditID,
		TenantID:  tenantID,
		Action:    "authentication.login",
		Resource:  "authentication",
		CreatedAt: time.Now().UTC(),
	}

	repository := &fakeAuditRepository{
		findResult: expectedAudit,
	}

	service := NewAuditService(repository)

	result, err := service.FindByID(
		context.Background(),
		tenantID,
		auditID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repository.findCalled {
		t.Fatal("expected repository FindByID to be called")
	}

	if repository.findTenantID != tenantID {
		t.Fatalf(
			"expected tenant ID %q, got %q",
			tenantID,
			repository.findTenantID,
		)
	}

	if repository.findAuditID != auditID {
		t.Fatalf(
			"expected audit ID %q, got %q",
			auditID,
			repository.findAuditID,
		)
	}

	if result != expectedAudit {
		t.Fatal("expected returned audit log to match repository result")
	}
}

func TestAuditServiceFindByIDValidation(t *testing.T) {
	tests := []struct {
		name        string
		tenantID    string
		auditID     string
		expectedErr error
	}{
		{
			name:        "missing tenant id",
			tenantID:    "",
			auditID:     "550e8400-e29b-41d4-a716-446655440000",
			expectedErr: ErrTenantIDRequired,
		},
		{
			name:        "missing audit id",
			tenantID:    "550e8400-e29b-41d4-a716-446655440000",
			auditID:     "",
			expectedErr: ErrAuditIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeAuditRepository{}
			service := NewAuditService(repository)

			_, err := service.FindByID(
				context.Background(),
				tt.tenantID,
				tt.auditID,
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedErr,
					err,
				)
			}

			if repository.findCalled {
				t.Fatal("repository should not be called for invalid input")
			}
		})
	}
}

func TestAuditServiceFindByIDInvalidUUID(t *testing.T) {
	repository := &fakeAuditRepository{}

	service := NewAuditService(repository)

	_, err := service.FindByID(
		context.Background(),
		"invalid-tenant-id",
		"550e8400-e29b-41d4-a716-446655440000",
	)

	if err == nil {
		t.Fatal("expected invalid tenant ID error")
	}

	if repository.findCalled {
		t.Fatal("repository should not be called with invalid tenant ID")
	}
}
