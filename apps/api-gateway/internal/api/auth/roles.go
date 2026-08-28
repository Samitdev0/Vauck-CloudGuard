package auth

import (
	"errors"
	"strings"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
	RoleUser     Role = "user"
)

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrForbidden   = errors.New("insufficient permissions")
)

func ParseRole(value string) (Role, error) {
	role := Role(strings.ToLower(strings.TrimSpace(value)))

	switch role {
	case RoleAdmin, RoleOperator, RoleViewer, RoleUser:
		return role, nil
	default:
		return "", ErrInvalidRole
	}
}

func HasRole(userRole string, allowed ...Role) bool {
	role, err := ParseRole(userRole)
	if err != nil {
		return false
	}

	for _, allowedRole := range allowed {
		if role == allowedRole {
			return true
		}
	}

	return false
}

func RequireRole(userRole string, allowed ...Role) error {
	if !HasRole(userRole, allowed...) {
		return ErrForbidden
	}

	return nil
}
