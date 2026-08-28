package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	query := `
		INSERT INTO users (
			tenant_id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.TenantID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.User, error) {

	query := `
		SELECT
			id,
			tenant_id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	tenantID string,
	email string,
) (*models.User, error) {

	query := `
		SELECT
			id,
			tenant_id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			active,
			created_at,
			updated_at
		FROM users
		WHERE tenant_id = $1
		  AND email = $2
	`

	var user models.User

	err := r.db.QueryRow(
		ctx,
		query,
		tenantID,
		email,
	).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user *models.User,
) error {

	query := `
		UPDATE users
		SET
			email = $1,
			password_hash = $2,
			first_name = $3,
			last_name = $4,
			role = $5,
			active = $6,
			updated_at = NOW()
		WHERE id = $7
		  AND tenant_id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		user.ID,
		user.TenantID,
	).Scan(
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(
	ctx context.Context,
	tenantID string,
	id string,
) error {

	query := `
		DELETE FROM users
		WHERE id = $1
		  AND tenant_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
		tenantID,
	)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
