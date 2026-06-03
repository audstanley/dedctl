package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server ServerConfig   `mapstructure:"server"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	Game   GameConfig     `mapstructure:"game"`
	Users  []UserConfig   `mapstructure:"users"`
}

// UserConfig holds a single user's credentials
type UserConfig struct {
	Username    string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
	IsAdmin     bool   `mapstructure:"is_admin"`
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

// LoadConfig loads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/steamctl/")
	viper.AddConfigPath("$HOME/.steamctl/")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.origins", []string{"http://localhost:5174"})
	viper.SetDefault("game.base_path", "$HOME/Games")
	viper.SetDefault("users", []UserConfig{
		{Username: "admin", PasswordHash: "", IsAdmin: true},
	})

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Parse config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".steamctl", "config.yaml")
}
