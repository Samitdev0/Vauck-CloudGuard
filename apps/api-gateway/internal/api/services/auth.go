package services

import (
	"context"
	"errors"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/models"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is inactive")
)

type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.JWT
}

func NewAuthService(
	users *repository.UserRepository,
	jwtService *auth.JWT,
) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwtService,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	tenantID string,
	email string,
	password string,
) (string, error) {

	user, err := s.users.FindByEmail(
		ctx,
		tenantID,
		email,
	)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	if !user.Active {
		return "", ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	token, err := s.jwt.GenerateToken(
		user.ID,
		user.TenantID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) CreateUser(
	ctx context.Context,
	user *models.User,
) error {

	if user.Role == "" {
		user.Role = "user"
	}

	user.Active = true

	return s.users.Create(ctx, user)
}
