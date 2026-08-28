package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTenantNotFound = errors.New("tenant not found")

type Tenant struct {
	ID   string
	Name string
	Slug string
}

type TenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (r *TenantRepository) Create(
	ctx context.Context,
	name string,
	slug string,
) (*Tenant, error) {
	const query = `
		INSERT INTO tenants (
			name,
			slug
		)
		VALUES ($1, $2)
		RETURNING id, name, slug
	`

	var tenant Tenant

	err := r.db.QueryRow(
		ctx,
		query,
		name,
		slug,
	).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
	)

	if err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	return &tenant, nil
}

func (r *TenantRepository) FindByID(
	ctx context.Context,
	id string,
) (*Tenant, error) {
	const query = `
		SELECT
			id,
			name,
			slug
		FROM tenants
		WHERE id = $1
	`

	var tenant Tenant

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find tenant by id: %w", err)
	}

	return &tenant, nil
}

func (r *TenantRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*Tenant, error) {
	const query = `
		SELECT
			id,
			name,
			slug
		FROM tenants
		WHERE slug = $1
	`

	var tenant Tenant

	err := r.db.QueryRow(
		ctx,
		query,
		slug,
	).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find tenant by slug: %w", err)
	}

	return &tenant, nil
}
