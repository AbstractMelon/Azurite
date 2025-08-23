package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Storage  StorageConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Path string
}

type AuthConfig struct {
	JWTSecret     string
	GitHubID      string
	GitHubSecret  string
	DiscordID     string
	DiscordSecret string
	GoogleID      string
	GoogleSecret  string
}

type StorageConfig struct {
	ModsPath    string
	ImagesPath  string
	MaxFileSize int64
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./azurite.db"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			GitHubID:      getEnv("GITHUB_CLIENT_ID", ""),
			GitHubSecret:  getEnv("GITHUB_CLIENT_SECRET", ""),
			DiscordID:     getEnv("DISCORD_CLIENT_ID", ""),
			DiscordSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
			GoogleID:      getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		Storage: StorageConfig{
			ModsPath:    getEnv("MODS_PATH", "./storage/mods"),
			ImagesPath:  getEnv("IMAGES_PATH", "./storage/images"),
			MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 104857600), // 100MB
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@azurite.dev"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
