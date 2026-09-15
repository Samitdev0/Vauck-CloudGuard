package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
)

var ErrAuditNotFound = errors.New("audit log not found")

type AuditRepository struct {
	db DBTX
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

func (r *AuditRepository) Create(
	ctx context.Context,
	audit *models.AuditLog,
) error {
	const query = `
		INSERT INTO audit_logs (
			id,
			tenant_id,
			user_id,
			action,
			resource,
			resource_id,
			ip_address,
			user_agent,
			metadata,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		audit.ID,
		audit.TenantID,
		audit.UserID,
		audit.Action,
		audit.Resource,
		audit.ResourceID,
		audit.IPAddress,
		audit.UserAgent,
		audit.Metadata,
		audit.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}

	return nil
}

func (r *AuditRepository) FindByID(
	ctx context.Context,
	tenantID string,
	id string,
) (*models.AuditLog, error) {
	const query = `
		SELECT
			id,
			tenant_id,
			user_id,
			action,
			resource,
			resource_id,
			ip_address,
			user_agent,
			metadata,
			created_at
		FROM audit_logs
		WHERE tenant_id = $1
		  AND id = $2
	`

	parsedTenantID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("parse tenant id: %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse audit id: %w", err)
	}

	var audit models.AuditLog

	err = r.db.QueryRow(
		ctx,
		query,
		parsedTenantID,
		parsedID,
	).Scan(
		&audit.ID,
		&audit.TenantID,
		&audit.UserID,
		&audit.Action,
		&audit.Resource,
		&audit.ResourceID,
		&audit.IPAddress,
		&audit.UserAgent,
		&audit.Metadata,
		&audit.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAuditNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find audit log: %w", err)
	}

	return &audit, nil
}
