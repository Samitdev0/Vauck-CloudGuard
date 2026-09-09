package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database/repository"
)

var (
	ErrAuditLogNil      = errors.New("audit log is nil")
	ErrTenantIDRequired = errors.New("tenant id is required")
	ErrAuditIDRequired  = errors.New("audit id is required")
	ErrActionRequired   = errors.New("audit action is required")
	ErrResourceRequired = errors.New("audit resource is required")
)

type AuditRepository interface {
	Create(
		context.Context,
		*models.AuditLog,
	) error

	FindByID(
		context.Context,
		string,
		string,
	) (*models.AuditLog, error)
}

type AuditService struct {
	repository AuditRepository
}

func NewAuditService(
	repository *repository.AuditRepository,
) *AuditService {
	return &AuditService{
		repository: repository,
	}
}

func NewAuditServiceWithRepository(
	repository AuditRepository,
) *AuditService {
	return &AuditService{
		repository: repository,
	}
}

func (s *AuditService) Create(
	ctx context.Context,
	audit *models.AuditLog,
) error {
	if audit == nil {
		return ErrAuditLogNil
	}

	if audit.TenantID == "" {
		return ErrTenantIDRequired
	}

	if audit.Action == "" {
		return ErrActionRequired
	}

	if audit.Resource == "" {
		return ErrResourceRequired
	}

	if audit.ID == "" {
		audit.ID = uuid.NewString()
	}

	if audit.CreatedAt.IsZero() {
		audit.CreatedAt = time.Now().UTC()
	}

	if err := s.repository.Create(ctx, audit); err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}

	return nil
}

func (s *AuditService) FindByID(
	ctx context.Context,
	tenantID string,
	id string,
) (*models.AuditLog, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if id == "" {
		return nil, ErrAuditIDRequired
	}

	if _, err := uuid.Parse(tenantID); err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}

	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid audit id: %w", err)
	}

	audit, err := s.repository.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("find audit log: %w", err)
	}

	return audit, nil
}
