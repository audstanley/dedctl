package service

import (
	"crypto/sha512"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"dedctl/internal/utils"
)

func hashPasswordSHA512(pw string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(pw)))
}

func hashPasswordBcrypt(pw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func TestLoginSuccessSHA512(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
		{Username: "operator", PasswordHash: hashPasswordSHA512("op456"), PasswordType: "sha512", IsAdmin: false},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	token, user, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
	if !user.IsAdmin {
		t.Error("expected admin user to have IsAdmin=true")
	}
}

func TestLoginSuccessSHA512DefaultType(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	token, user, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
}

func TestLoginSuccessBcrypt(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordBcrypt("admin123"), PasswordType: "bcrypt", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	token, user, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
}

func TestLoginSuccessPlain(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: "admin123", PasswordType: "plain", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	token, user, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
}

func TestLoginSuccessOperator(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
		{Username: "operator", PasswordHash: hashPasswordSHA512("op456"), PasswordType: "sha512", IsAdmin: false},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, user, err := auth.Login("operator", "op456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.IsAdmin {
		t.Error("expected operator user to have IsAdmin=false")
	}
}

func TestLoginWrongPasswordSHA512(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("admin", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLoginWrongPasswordBcrypt(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordBcrypt("admin123"), PasswordType: "bcrypt", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("admin", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLoginWrongPasswordPlain(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: "admin123", PasswordType: "plain", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("admin", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLoginWrongUsername(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("nobody", "admin123")
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
}

func TestLoginEmptyPassword(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512(""), PasswordType: "sha512", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("admin", "")
	if err != nil {
		t.Fatalf("expected no error for empty password match, got %v", err)
	}
}

func TestLoginInvalidPasswordType(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: "somehash", PasswordType: "invalid_type", IsAdmin: true},
	}
	secret := "test-jwt-secret"

	auth := NewAuthService(users, secret)

	_, _, err := auth.Login("admin", "anypassword")
	if err == nil {
		t.Fatal("expected error for invalid password type, got nil")
	}
}

func TestTokenIsValidAgainstUtils(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
	}
	secret := "my-jwt-secret"

	auth := NewAuthService(users, secret)

	token, _, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	claims, err := utils.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("token validation failed: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", claims.Username)
	}
}

func TestTokenFailsWithWrongSecret(t *testing.T) {
	users := []UserInfo{
		{Username: "admin", PasswordHash: hashPasswordSHA512("admin123"), PasswordType: "sha512", IsAdmin: true},
	}
	secret := "my-jwt-secret"

	auth := NewAuthService(users, secret)

	token, _, err := auth.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	_, err = utils.ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected validation error with wrong secret, got nil")
	}
}

func TestEmptyUsers(t *testing.T) {
	auth := NewAuthService([]UserInfo{}, "secret")

	_, _, err := auth.Login("admin", "admin123")
	if err == nil {
		t.Fatal("expected error for empty user list, got nil")
	}
}

func TestLoginReturnsValidJWTFormat(t *testing.T) {
	users := []UserInfo{
		{Username: "testuser", PasswordHash: hashPasswordSHA512("testpass"), PasswordType: "sha512", IsAdmin: false},
	}
	secret := "test-secret"

	auth := NewAuthService(users, secret)

	token, _, err := auth.Login("testuser", "testpass")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("expected JWT token with 3 parts, got %d", len(parts))
	}
}
