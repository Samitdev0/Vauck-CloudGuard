package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultTenantName = "CloudGuard Local"
	DefaultTenantSlug = "cloudguard-local"

	DefaultAdminEmail    = "admin@cloudguard.local"
	DefaultAdminPassword = "CloudGuard123!"
)

func Run(
	ctx context.Context,
	tenantRepository *repository.TenantRepository,
	userRepository *repository.UserRepository,
) error {

	tenant, err := tenantRepository.FindBySlug(
		ctx,
		DefaultTenantSlug,
	)

	if err != nil {
		if !errors.Is(err, repository.ErrTenantNotFound) {
			return fmt.Errorf("find default tenant: %w", err)
		}

		tenant, err = tenantRepository.Create(
			ctx,
			DefaultTenantName,
			DefaultTenantSlug,
		)
		if err != nil {
			return fmt.Errorf("create default tenant: %w", err)
		}
	}

	_, err = userRepository.FindByEmail(
		ctx,
		tenant.ID,
		DefaultAdminEmail,
	)

	if err == nil {
		return nil
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("find default admin: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(DefaultAdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("hash default admin password: %w", err)
	}

	admin := &models.User{
		TenantID:     tenant.ID,
		Email:        DefaultAdminEmail,
		PasswordHash: string(passwordHash),
		FirstName:    "CloudGuard",
		LastName:     "Administrator",
		Role:         "admin",
		Active:       true,
	}

	if err := userRepository.Create(ctx, admin); err != nil {
		return fmt.Errorf("create default admin: %w", err)
	}

	return nil
}
