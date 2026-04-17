package repository

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestFindByUsername_Found(t *testing.T) {
	repo := NewUserRepository()

	user, err := repo.FindByUsername("admin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", user.Role)
	}
	if user.ID == 0 {
		t.Error("expected non-zero user ID")
	}
}

func TestFindByUsername_NotFound(t *testing.T) {
	repo := NewUserRepository()

	_, err := repo.FindByUsername("nonexistent")
	if err == nil {
		t.Error("expected error for unknown username, got nil")
	}
}

func TestDefaultAccounts_AtLeastTwo(t *testing.T) {
	repo := NewUserRepository()

	usernames := []string{"admin", "user"}
	for _, name := range usernames {
		u, err := repo.FindByUsername(name)
		if err != nil {
			t.Errorf("expected default account '%s' to exist, got error: %v", name, err)
			continue
		}
		if u.Username != name {
			t.Errorf("expected username '%s', got '%s'", name, u.Username)
		}
	}
}

func TestDefaultAccounts_PasswordIsHashed(t *testing.T) {
	repo := NewUserRepository()

	cases := []struct{ username, plaintext string }{
		{"admin", "admin123"},
		{"user", "user123"},
	}
	for _, tc := range cases {
		u, err := repo.FindByUsername(tc.username)
		if err != nil {
			t.Fatalf("user '%s' not found: %v", tc.username, err)
		}
		if u.PasswordHash == tc.plaintext {
			t.Errorf("password for '%s' should not be stored as plaintext", tc.username)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(tc.plaintext)); err != nil {
			t.Errorf("bcrypt verification failed for '%s': %v", tc.username, err)
		}
	}
}

func TestFindByUsername_ReturnsCopy(t *testing.T) {
	repo := NewUserRepository()

	u1, _ := repo.FindByUsername("admin")
	u1.Username = "modified"

	u2, _ := repo.FindByUsername("admin")
	if u2.Username != "admin" {
		t.Error("FindByUsername should return a copy, not a reference")
	}
}
