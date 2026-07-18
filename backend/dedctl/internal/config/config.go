package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// viperInstance is a package-level viper instance to allow test overrides
var viperInstance = viper.New()

// Config holds the application configuration
type Config struct {
	Server ServerConfig   `mapstructure:"server"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	Game   GameConfig     `mapstructure:"game"`
	Users  []UserConfig   `mapstructure:"users"`
}

// UserConfig holds a single user's credentials
type UserConfig struct {
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
	PasswordType string `mapstructure:"password_type"`
	IsAdmin      bool   `mapstructure:"is_admin"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    string   `mapstructure:"port"`
	Host    string   `mapstructure:"host"`
	Origins []string `mapstructure:"origins"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	ExpiresIn string `mapstructure:"expires_in"`
}

// GameConfig holds game server configuration
type GameConfig struct {
	BasePath string `mapstructure:"base_path"`
}

var externalConfigFile string

// SetConfigFile sets an external config file path to use
func SetConfigFile(path string) {
	externalConfigFile = path
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	hasFile := viperInstance.ConfigFileUsed() != ""

	if !hasFile {
		if externalConfigFile != "" {
			viperInstance.SetConfigFile(externalConfigFile)
		} else {
			viperInstance.SetConfigName("config")
			viperInstance.SetConfigType("yaml")
			home, _ := os.UserHomeDir()
			viperInstance.AddConfigPath(filepath.Join(home, ".dedctl"))
			viperInstance.AddConfigPath("/etc/dedctl/")
		}
	}

	// Set default values
	viperInstance.SetDefault("server.port", "8080")
	viperInstance.SetDefault("server.host", "0.0.0.0")
	viperInstance.SetDefault("server.origins", []string{"http://localhost:5174"})
	viperInstance.SetDefault("game.base_path", "$HOME/Games")
	viperInstance.SetDefault("users", []UserConfig{
		{Username: "admin", PasswordHash: "", IsAdmin: true},
	})

	// Read config
	if err := viperInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	// Read environment variables
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Parse config
	var config Config
	if err := viperInstance.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dedctl", "config.yaml")
}

// GetConfigFileUsed returns the path to the config file that was used, or empty string.
func GetConfigFileUsed() string {
	return viperInstance.ConfigFileUsed()
}
