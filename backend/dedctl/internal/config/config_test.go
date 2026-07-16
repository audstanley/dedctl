package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func newTestViper() *viper.Viper {
	v := viper.New()
	return v
}

func TestLoadConfigDefaults(t *testing.T) {
	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigName("nonexistent")
	viperInstance.SetConfigType("yaml")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Server.Port != "8085" {
		t.Errorf("expected default port '8085', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `server:
  port: "9090"
  host: "127.0.0.1"
  origins:
    - "http://example.com"
jwt:
  secret_key: "my-secret"
  expires_in: "12h"
game:
  base_path: "/opt/games"
users:
  - username: testuser
    password_hash: "abc123"
    password_type: sha512
    is_admin: true
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected port '9090', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host '127.0.0.1', got '%s'", cfg.Server.Host)
	}
	if cfg.JWT.SecretKey != "my-secret" {
		t.Errorf("expected secret 'my-secret', got '%s'", cfg.JWT.SecretKey)
	}
	if cfg.JWT.ExpiresIn != "12h" {
		t.Errorf("expected expires_in '12h', got '%s'", cfg.JWT.ExpiresIn)
	}
	if cfg.Game.BasePath != "/opt/games" {
		t.Errorf("expected base_path '/opt/games', got '%s'", cfg.Game.BasePath)
	}
	if len(cfg.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(cfg.Users))
	}
	if cfg.Users[0].Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", cfg.Users[0].Username)
	}
	if cfg.Users[0].PasswordHash != "abc123" {
		t.Errorf("expected password_hash 'abc123', got '%s'", cfg.Users[0].PasswordHash)
	}
	if cfg.Users[0].PasswordType != "sha512" {
		t.Errorf("expected password_type 'sha512', got '%s'", cfg.Users[0].PasswordType)
	}
	if !cfg.Users[0].IsAdmin {
		t.Error("expected is_admin=true")
	}
}

func TestLoadConfigMissingField(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `server:
  port: "8888"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Server.Port != "8888" {
		t.Errorf("expected port '8888', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `server:
  port: "8888"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() {
		viperInstance = originalViper
		os.Unsetenv("SERVER_PORT")
	}()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)
	os.Setenv("SERVER_PORT", "7777")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Server.Port != "7777" {
		t.Errorf("expected port overridden to '7777', got '%s'", cfg.Server.Port)
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Fatal("expected non-empty config path")
	}
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got '%s'", path)
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected config file name 'config.yaml', got '%s'", filepath.Base(path))
	}
}

func TestLoadConfigUserWithPasswordTypeBcrypt(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `users:
  - username: bcrypt_user
    password_hash: "$2b$12$example"
    password_type: bcrypt
    is_admin: false
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cfg.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(cfg.Users))
	}
	if cfg.Users[0].PasswordType != "bcrypt" {
		t.Errorf("expected password_type 'bcrypt', got '%s'", cfg.Users[0].PasswordType)
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `server:
  port: "8888"
  invalid: [yaml: {{{`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("expected yaml error, got: %v", err)
	}
}

func TestLoadConfigUserWithPasswordTypePlain(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `users:
  - username: plain_user
    password_hash: "mysecret"
    password_type: plain
    is_admin: false
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalViper := viperInstance
	defer func() { viperInstance = originalViper }()

	viperInstance = newTestViper()
	viperInstance.SetConfigFile(configPath)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cfg.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(cfg.Users))
	}
	if cfg.Users[0].PasswordType != "plain" {
		t.Errorf("expected password_type 'plain', got '%s'", cfg.Users[0].PasswordType)
	}
}
