package auth

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrIdentityNotFound = errors.New("authenticated identity not found")
)

func GetClaims(c echo.Context) (*Claims, error) {
	value := c.Get("user")

	if value == nil {
		return nil, ErrIdentityNotFound
	}

	claims, ok := value.(*Claims)
	if !ok || claims == nil {
		return nil, ErrIdentityNotFound
	}

	return claims, nil
}

func GetUserID(c echo.Context) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}

	if claims.UserID == "" {
		return "", ErrIdentityNotFound
	}

	return claims.UserID, nil
}

func GetTenantID(c echo.Context) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}

	if claims.TenantID == "" {
		return "", ErrIdentityNotFound
	}

	return claims.TenantID, nil
}

func GetRole(c echo.Context) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}

	if claims.Role == "" {
		return "", ErrIdentityNotFound
	}

	return claims.Role, nil
}

func GetEmail(c echo.Context) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}

	if claims.Email == "" {
		return "", ErrIdentityNotFound
	}

	return claims.Email, nil
}
