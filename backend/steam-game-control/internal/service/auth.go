package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

// UserService handles user operations
type UserService struct {
	db *leveldb.DB
}

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

// NewUserService creates a new UserService
func NewUserService() (*UserService, error) {
	// Create data directory if it doesn't exist
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	db, err := leveldb.OpenFile(filepath.Join(dataDir, "users"), nil)
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

// Close closes the database connection
func (s *UserService) Close() {
	s.db.Close()
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, password string, isAdmin bool) error {
	user := User{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
	}

	// Store user data
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.db.Put([]byte(username), userData, nil)
}

// GetUser retrieves a user by username
func (s *UserService) GetUser(username string) (*User, error) {
	data, err := s.db.Get([]byte(username), nil)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Authenticate validates user credentials
func (s *UserService) Authenticate(username, password string) (*User, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}

	// In a real implementation, you would compare hashed passwords
	// For this demo, we're doing a simple string comparison
	if user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// AuthService handles authentication operations
type AuthService struct {
	userService *UserService
}

// NewAuthService creates a new AuthService
func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

// Login authenticates a user
func (s *AuthService) Login(username, password string) (string, *User, error) {
	// In a real implementation, this would validate credentials
	// and return a JWT token

	user, err := s.userService.Authenticate(username, password)
	if err != nil {
		return "", nil, err
	}

	// For demo purposes, just returning a dummy token
	// In reality, this should generate a real JWT token
	return "dummy-jwt-token", user, nil
}

// Register creates a new user
func (s *AuthService) Register(username, password string, isAdmin bool) error {
	// Check if user already exists
	_, err := s.userService.GetUser(username)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	return s.userService.CreateUser(username, password, isAdmin)
}
