package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/config"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrMissingJWTSecret = errors.New("jwt secret is not configured")
)

type Claims struct {
	UserID   string `json:"userId"`
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
	issuer string
}

func NewJWT(cfg config.JWTConfig) (*JWT, error) {
	secret := strings.TrimSpace(cfg.Secret)

	if secret == "" {
		return nil, ErrMissingJWTSecret
	}

	if len(secret) < 32 {
		return nil, errors.New(
			"jwt secret must contain at least 32 characters",
		)
	}

	issuer := strings.TrimSpace(cfg.Issuer)
	if issuer == "" {
		issuer = "CloudGuard"
	}

	return &JWT{
		secret: []byte(secret),
		issuer: issuer,
	}, nil
}

func (j *JWT) GenerateToken(
	userID string,
	tenantID string,
	email string,
	role string,
) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"%w: unexpected signing method %s",
					ErrInvalidToken,
					token.Method.Alg(),
				)
			}

			return j.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuer(j.issuer),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"%w: %v",
			ErrInvalidToken,
			err,
		)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.UserID == "" ||
		claims.TenantID == "" ||
		claims.Email == "" ||
		claims.Role == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
