package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Name      string `yaml:"name"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse_time"`
	Loc       string `yaml:"loc"`
}

type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn int    `yaml:"expires_in"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var appConfig *Config

func LoadConfig(configPath string) error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional
	}

	// Read YAML config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &appConfig); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Override with environment variables
	overrideWithEnvVars()

	return nil
}

func overrideWithEnvVars() {
	// App config
	if port := os.Getenv("APP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			appConfig.App.Port = p
		}
	}
	if mode := os.Getenv("APP_MODE"); mode != "" {
		appConfig.App.Mode = mode
	}

	// Database config
	if host := os.Getenv("DB_HOST"); host != "" {
		appConfig.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			appConfig.Database.Port = p
		}
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		appConfig.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		appConfig.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		appConfig.Database.Name = name
	}

	// JWT config
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		appConfig.JWT.Secret = secret
	}
	if expiresIn := os.Getenv("JWT_EXPIRES_IN"); expiresIn != "" {
		if e, err := strconv.Atoi(expiresIn); err == nil {
			appConfig.JWT.ExpiresIn = e
		}
	}

	// Redis config
	if host := os.Getenv("REDIS_HOST"); host != "" {
		appConfig.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			appConfig.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		appConfig.Redis.Password = password
	}

	// Logger config
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		appConfig.Logger.Level = level
	}
}

func GetConfig() *Config {
	return appConfig
}
