package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUsername     string
	DBPassword     string
	DBName         string
	DBHost         string
	Port           string
	ExtensionID    string
	DevClientURL   string
	AllowedOrigins string
	Environment    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	config := &Config{
		DBUsername:     getEnv("DB_USERNAME", ""),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", ""),
		DBHost:         getEnv("DB_HOST", ""),
		Port:           getEnv("PORT", ""),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", ""),
		Environment:    getEnv("ENVIRONMENT", "RELEASE"),
	}

	return config
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		c.DBUsername,
		c.DBPassword,
		c.DBHost,
		c.DBName,
	)
}

func (c *Config) Validate() error {
	if c.DBUsername == "" || c.DBPassword == "" || c.DBName == "" {
		return fmt.Errorf("데이터베이스 구성이 완료되지 않았습니다")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
