package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	LogLevel string `yaml:"log_level"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SessionConfig struct {
	Name     string `yaml:"name"`
	Secret   string `yaml:"secret"`
	Lifetime int    `yaml:"lifetime"` // in minutes
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Session  SessionConfig  `yaml:"session"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Driver:   "sqlite3",
			Name:     "django.db",
			Path:     "./data",
			LogLevel: "info",
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Session: SessionConfig{
			Name:     "django_session",
			Secret:   "change-this-secret-key",
			Lifetime: 120, // 2 hours
		},
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the YAML into our Config struct
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Ensure the database directory exists
	if err := os.MkdirAll(config.Database.Path, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	return config, nil
}

// SaveConfig saves the current configuration to a YAML file
func (c *Config) Save(path string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Marshal the config to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
