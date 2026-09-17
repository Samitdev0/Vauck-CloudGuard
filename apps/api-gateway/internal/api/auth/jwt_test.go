package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/config"
)

func testJWT(t *testing.T) *JWT {
	t.Helper()

	j, err := NewJWT(config.JWTConfig{
		Secret: "cloudguard-test-secret-0123456789-abcd",
		Issuer: "CloudGuard",
	})
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	return j
}

func TestJWTValidateTokenReturnsClaims(t *testing.T) {
	j := testJWT(t)

	token, err := j.GenerateToken(
		"user-a",
		"tenant-a",
		"user@example.com",
		"admin",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := j.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}

	if claims.UserID != "user-a" {
		t.Fatalf("expected user ID %q, got %q", "user-a", claims.UserID)
	}

	if claims.TenantID != "tenant-a" {
		t.Fatalf("expected tenant ID %q, got %q", "tenant-a", claims.TenantID)
	}

	if claims.Email != "user@example.com" {
		t.Fatalf("expected email %q, got %q", "user@example.com", claims.Email)
	}

	if claims.Role != "admin" {
		t.Fatalf("expected role %q, got %q", "admin", claims.Role)
	}
}

func TestJWTValidateTokenRejectsTamperedTenant(t *testing.T) {
	j := testJWT(t)

	token, err := j.GenerateToken(
		"user-a",
		"tenant-a",
		"user@example.com",
		"admin",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected JWT with 3 parts, got %d", len(parts))
	}

	tamperedClaims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			UserID:   "user-a",
			TenantID: "tenant-b",
			Email:    "user@example.com",
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   "user-a",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		},
	)

	tamperedParts := strings.Split(
		mustSignedTokenWithoutSecret(t, tamperedClaims),
		".",
	)

	if len(tamperedParts) != 3 {
		t.Fatalf("expected tampered JWT with 3 parts, got %d", len(tamperedParts))
	}

	tamperedToken := parts[0] + "." + tamperedParts[1] + "." + parts[2]

	_, err = j.ValidateToken(tamperedToken)
	if err == nil {
		t.Fatal("expected tampered token to be rejected")
	}

	if !strings.Contains(err.Error(), ErrInvalidToken.Error()) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTValidateTokenRejectsMissingTenant(t *testing.T) {
	j := testJWT(t)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			UserID:   "user-a",
			TenantID: "",
			Email:    "user@example.com",
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   "user-a",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		},
	)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = j.ValidateToken(signedToken)
	if err == nil {
		t.Fatal("expected token without tenant to be rejected")
	}

	if !strings.Contains(err.Error(), ErrInvalidToken.Error()) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTValidateTokenRejectsExpiredToken(t *testing.T) {
	j := testJWT(t)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			UserID:   "user-a",
			TenantID: "tenant-a",
			Email:    "user@example.com",
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   "user-a",
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-30 * time.Minute)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-15 * time.Minute)),
			},
		},
	)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = j.ValidateToken(signedToken)
	if err == nil {
		t.Fatal("expected expired token to be rejected")
	}

	if !strings.Contains(err.Error(), ErrInvalidToken.Error()) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTValidateTokenRejectsUnexpectedAlgorithm(t *testing.T) {
	j := testJWT(t)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS384,
		Claims{
			UserID:   "user-a",
			TenantID: "tenant-a",
			Email:    "user@example.com",
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				Subject:   "user-a",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		},
	)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = j.ValidateToken(signedToken)
	if err == nil {
		t.Fatal("expected token with unexpected algorithm to be rejected")
	}

	if !strings.Contains(err.Error(), ErrInvalidToken.Error()) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func mustSignedTokenWithoutSecret(t *testing.T, token *jwt.Token) string {
	t.Helper()

	// Generate a token using a deliberately different secret.
	// The resulting signature is intentionally invalid for the CloudGuard JWT.
	signedToken, err := token.SignedString(
		[]byte("attacker-controlled-secret-0123456789"),
	)
	if err != nil {
		t.Fatalf("failed to sign tampered token: %v", err)
	}

	return signedToken
}
