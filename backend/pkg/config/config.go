package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Auth    AuthConfig    `mapstructure:"auth"`
	Agent   AgentConfig   `mapstructure:"agent"`
	Storage StorageConfig `mapstructure:"storage"`
	Rustun  RustunConfig  `mapstructure:"rustun"` // Legacy, for backward compatibility
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type AuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type AgentConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Provider       string `mapstructure:"provider"`         // openai, deepseek, ollama
	OpenAIAPIKey   string `mapstructure:"openai_api_key"`   // For OpenAI
	DeepSeekAPIKey string `mapstructure:"deepseek_api_key"` // For DeepSeek
	Model          string `mapstructure:"model"`
	BaseURL        string `mapstructure:"base_url"`        // Custom API endpoint
	LocalModelURL  string `mapstructure:"local_model_url"` // For future Ollama support
}

type StorageConfig struct {
	Type     string         `mapstructure:"type"` // "file" or "database"
	File     FileConfig     `mapstructure:"file"`
	Database DatabaseConfig `mapstructure:"database"`
}

type FileConfig struct {
	RoutesFile         string `mapstructure:"routes_file"`
	RoutesFileFallback string `mapstructure:"routes_file_fallback"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // mysql, postgres, sqlite
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Path     string `mapstructure:"path"` // For SQLite
}

type RustunConfig struct {
	RoutesFile         string `mapstructure:"routes_file"`
	RoutesFileFallback string `mapstructure:"routes_file_fallback"`
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("auth.username", "admin")
	v.SetDefault("auth.password", "admin123")

	v.SetDefault("storage.type", "file")
	v.SetDefault("storage.file.routes_file", "/etc/rustun/routes.json")
	v.SetDefault("storage.file.routes_file_fallback", "./routes.json")

	v.SetDefault("storage.database.type", "mysql")
	v.SetDefault("storage.database.host", "localhost")
	v.SetDefault("storage.database.port", 3306)

	v.SetDefault("agent.enabled", true)
	v.SetDefault("agent.provider", "openai")
	v.SetDefault("agent.model", "gpt-4o-mini")

	// Legacy defaults
	v.SetDefault("rustun.routes_file", "/etc/rustun/routes.json")
	v.SetDefault("rustun.routes_file_fallback", "./routes.json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Backward compatibility: use rustun config if storage.file is not set
	if config.Storage.File.RoutesFile == "" {
		config.Storage.File.RoutesFile = config.Rustun.RoutesFile
		config.Storage.File.RoutesFileFallback = config.Rustun.RoutesFileFallback
	}

	// Determine which routes file to use (for file storage)
	if config.Storage.Type == "file" {
		if _, err := os.Stat(config.Storage.File.RoutesFile); os.IsNotExist(err) {
			config.Storage.File.RoutesFile = config.Storage.File.RoutesFileFallback
		}
	}

	return &config, nil
}

// Address returns the full server address
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DSN returns the database connection string for MySQL
func (c *DatabaseConfig) DSN() string {
	switch c.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			c.Host, c.Username, c.Password, c.Database, c.Port)
	case "sqlite":
		if c.Path != "" {
			return c.Path
		}
		return "./rustun.db"
	default:
		return ""
	}
}
