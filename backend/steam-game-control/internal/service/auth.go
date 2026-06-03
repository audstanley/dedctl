package service

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
	PasswordType string
	IsAdmin      bool
}

// NewAuthService creates a new AuthService with users and JWT secret
func NewAuthService(users []UserInfo, secretKey string) *AuthService {
	return &AuthService{
		users:     users,
		secretKey: secretKey,
	}
}

// verifyPassword compares the provided password against the stored hash using the specified algorithm
func verifyPassword(password, storedHash, passwordType string) bool {
	switch passwordType {
	case "bcrypt":
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
		return err == nil
	case "plain":
		return password == storedHash
	case "":
		fallthrough
	case "sha512":
		hash := fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
		return hash == storedHash
	default:
		return false
	}
}

// Login authenticates a user by comparing the password against stored hashes
func (s *AuthService) Login(username, password string) (string, *User, error) {
	for _, u := range s.users {
		if u.Username == username && verifyPassword(password, u.PasswordHash, u.PasswordType) {
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
