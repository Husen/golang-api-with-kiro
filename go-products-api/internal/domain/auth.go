package domain

import "github.com/golang-jwt/jwt/v5"

// User represents the user entity
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string // "admin" | "user"
}

// LoginRequest is the payload for the login endpoint
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response after successful login
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// JWTClaims holds the JWT claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// UserRepository defines the contract for user data access
type UserRepository interface {
	FindByUsername(username string) (*User, error)
}

// AuthUsecase defines the contract for authentication business logic
type AuthUsecase interface {
	Login(username, password string) (string, error)
}
