package repository

import (
	"errors"
	"go-products-api/internal/domain"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	mu    sync.Mutex
	users []domain.User
}

func NewUserRepository() domain.UserRepository {
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

	return &userRepository{
		users: []domain.User{
			{ID: 1, Username: "admin", PasswordHash: string(adminHash), Role: "admin"},
			{ID: 2, Username: "user", PasswordHash: string(userHash), Role: "user"},
		},
	}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Username == username {
			cp := u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}
