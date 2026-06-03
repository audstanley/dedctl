package service

import (
	"crypto/sha256"
	"fmt"

	"steam-game-control/internal/utils"
)

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

// AuthService handles authentication operations
type AuthService struct {
	users     []UserInfo
	secretKey string
}

// UserInfo holds user configuration data
type UserInfo struct {
	Username     string
	PasswordHash string
	IsAdmin      bool
}

// NewAuthService creates a new AuthService with users and JWT secret
func NewAuthService(users []UserInfo, secretKey string) *AuthService {
	return &AuthService{
		users:     users,
		secretKey: secretKey,
	}
}

// Login authenticates a user by comparing SHA-256 hash of the password
func (s *AuthService) Login(username, password string) (string, *User, error) {
	inputHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	for _, u := range s.users {
		if u.Username == username && u.PasswordHash == inputHash {
			token, err := utils.GenerateToken(username, s.secretKey)
			if err != nil {
				return "", nil, fmt.Errorf("failed to generate token: %v", err)
			}

			return token, &User{
				Username: u.Username,
				IsAdmin:  u.IsAdmin,
			}, nil
		}
	}

	return "", nil, fmt.Errorf("invalid credentials")
}
