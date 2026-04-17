package usecase

import (
	"os"
	"testing"
	"time"

	"go-products-api/internal/domain"
	"go-products-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

func newTestUsecase(t *testing.T) domain.AuthUsecase {
	t.Helper()
	repo := repository.NewUserRepository()
	return NewAuthUsecase(repo)
}

func TestLogin_ValidCredentials(t *testing.T) {
	cases := []struct{ username, password string }{
		{"admin", "admin123"},
		{"user", "user123"},
	}
	for _, tc := range cases {
		uc := newTestUsecase(t)
		token, err := uc.Login(tc.username, tc.password)
		if err != nil {
			t.Errorf("expected no error for %s, got: %v", tc.username, err)
		}
		if token == "" {
			t.Errorf("expected non-empty token for %s", tc.username)
		}
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	uc := newTestUsecase(t)
	_, err := uc.Login("admin", "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestLogin_UsernameNotFound(t *testing.T) {
	uc := newTestUsecase(t)
	_, err := uc.Login("nonexistent", "password")
	if err == nil {
		t.Error("expected error for unknown username, got nil")
	}
}

func TestLogin_JWTSecretFromEnv(t *testing.T) {
	secret := "test-env-secret"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	repo := repository.NewUserRepository()
	uc := NewAuthUsecase(repo)

	tokenStr, err := uc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Parse and verify using the env secret
	claims := &domain.JWTClaims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		t.Errorf("token not valid with env secret: %v", err)
	}
}

func TestLogin_FallbackDefaultSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	repo := repository.NewUserRepository()
	uc := NewAuthUsecase(repo)

	tokenStr, err := uc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims := &domain.JWTClaims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("dev-secret-key"), nil
	})
	if err != nil || !tok.Valid {
		t.Errorf("token not valid with default secret: %v", err)
	}
}

func TestLogin_ClaimsContent(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	repo := repository.NewUserRepository()
	uc := NewAuthUsecase(repo)

	tokenStr, err := uc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims := &domain.JWTClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("dev-secret-key"), nil
	})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if claims.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", claims.Username)
	}
	if claims.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", claims.Role)
	}
	if claims.UserID == 0 {
		t.Error("expected non-zero user_id")
	}
	// exp should be ~24h from now
	exp := claims.ExpiresAt.Time
	diff := exp.Sub(time.Now())
	if diff < 23*time.Hour || diff > 25*time.Hour {
		t.Errorf("expected expiry ~24h from now, got %v", diff)
	}
}
