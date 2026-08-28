package auth

import (
	"github.com/labstack/echo/v4"
)

const ContextIdentity = "identity"

type Identity struct {
	UserID   string
	TenantID string
	Email    string
	Role     string
}

func IdentityFromClaims(claims *Claims) (*Identity, error) {
	if claims == nil {
		return nil, ErrIdentityNotFound
	}

	if claims.UserID == "" ||
		claims.TenantID == "" ||
		claims.Email == "" ||
		claims.Role == "" {
		return nil, ErrIdentityNotFound
	}

	return &Identity{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Email:    claims.Email,
		Role:     claims.Role,
	}, nil
}

func GetIdentity(c echo.Context) (*Identity, error) {
	value := c.Get(ContextIdentity)

	identity, ok := value.(*Identity)
	if !ok || identity == nil {
		return nil, ErrIdentityNotFound
	}

	return identity, nil
}
