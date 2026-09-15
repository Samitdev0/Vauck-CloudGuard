package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestAuditRepositoryFindByIDRequiresTenantIsolation(t *testing.T) {
	tenantA := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	auditB := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	db := &mockDB{
		row: &mockRow{
			scanFunc: func(dest ...any) error {
				return pgx.ErrNoRows
			},
		},
	}

	repository := &AuditRepository{
		db: db,
	}

	_, err := repository.FindByID(
		context.Background(),
		tenantA.String(),
		auditB.String(),
	)

	if !errors.Is(err, ErrAuditNotFound) {
		t.Fatalf(
			"expected ErrAuditNotFound, got %v",
			err,
		)
	}

	if len(db.args) != 2 {
		t.Fatalf(
			"expected 2 query arguments, got %d",
			len(db.args),
		)
	}

	if db.args[0] != tenantA {
		t.Fatalf(
			"expected first argument to be tenant ID %v, got %v",
			tenantA,
			db.args[0],
		)
	}

	if db.args[1] != auditB {
		t.Fatalf(
			"expected second argument to be audit ID %v, got %v",
			auditB,
			db.args[1],
		)
	}
}

func TestAuditRepositoryFindByIDReturnsAuditForMatchingTenant(t *testing.T) {
	tenantA := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	auditA := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	db := &mockDB{
		row: &mockRow{
			scanFunc: func(dest ...any) error {
				*dest[0].(*string) = auditA.String()
				*dest[1].(*string) = tenantA.String()

				return nil
			},
		},
	}

	repository := &AuditRepository{
		db: db,
	}

	audit, err := repository.FindByID(
		context.Background(),
		tenantA.String(),
		auditA.String(),
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if audit == nil {
		t.Fatal("expected audit log, got nil")
	}

	if audit.ID != auditA.String() {
		t.Fatalf(
			"expected audit ID %q, got %q",
			auditA.String(),
			audit.ID,
		)
	}

	if audit.TenantID != tenantA.String() {
		t.Fatalf(
			"expected tenant ID %q, got %q",
			tenantA.String(),
			audit.TenantID,
		)
	}
}
