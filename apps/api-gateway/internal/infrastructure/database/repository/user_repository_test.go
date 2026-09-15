package repository

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockRow struct {
	scanFunc func(dest ...any) error
}

func (r *mockRow) Scan(dest ...any) error {
	return r.scanFunc(dest...)
}

type mockDB struct {
	query string
	args  []any
	row   pgx.Row
}

func (m *mockDB) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) pgx.Row {
	m.query = query
	m.args = args

	return m.row
}

func (m *mockDB) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (pgconn.CommandTag, error) {
	panic("not implemented")
}

func TestUserRepositoryFindByIDRequiresTenantIsolation(t *testing.T) {
	const (
		tenantA = "11111111-1111-1111-1111-111111111111"
		userB   = "22222222-2222-2222-2222-222222222222"
	)

	db := &mockDB{
		row: &mockRow{
			scanFunc: func(dest ...any) error {
				return pgx.ErrNoRows
			},
		},
	}

	repository := &UserRepository{
		db: db,
	}

	_, err := repository.FindByID(
		context.Background(),
		tenantA,
		userB,
	)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}

	if !strings.Contains(db.query, "tenant_id = $2") {
		t.Fatalf(
			"expected query to filter by tenant_id, got: %s",
			db.query,
		)
	}

	if len(db.args) != 2 {
		t.Fatalf(
			"expected 2 query arguments, got %d",
			len(db.args),
		)
	}

	if db.args[0] != userB {
		t.Fatalf(
			"expected first argument to be user ID %q, got %v",
			userB,
			db.args[0],
		)
	}

	if db.args[1] != tenantA {
		t.Fatalf(
			"expected second argument to be tenant ID %q, got %v",
			tenantA,
			db.args[1],
		)
	}
}

func TestUserRepositoryFindByIDReturnsUserForMatchingTenant(t *testing.T) {
	const (
		tenantA = "11111111-1111-1111-1111-111111111111"
		userA   = "22222222-2222-2222-2222-222222222222"
	)

	db := &mockDB{
		row: &mockRow{
			scanFunc: func(dest ...any) error {
				values := []any{
					userA,
					tenantA,
					"user@example.com",
					"hashed-password",
					"John",
					"Doe",
					"admin",
					true,
				}

				for i := range values {
					switch destination := dest[i].(type) {
					case *string:
						*destination = values[i].(string)
					case *bool:
						*destination = values[i].(bool)
					}
				}

				return nil
			},
		},
	}

	repository := &UserRepository{
		db: db,
	}

	user, err := repository.FindByID(
		context.Background(),
		tenantA,
		userA,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != userA {
		t.Fatalf(
			"expected user ID %q, got %q",
			userA,
			user.ID,
		)
	}

	if user.TenantID != tenantA {
		t.Fatalf(
			"expected tenant ID %q, got %q",
			tenantA,
			user.TenantID,
		)
	}
}
