package config

import "os"

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	Port          string
	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "ottuser"),
		DBPassword:    getEnv("DB_PASSWORD", "ottpass123"),
		DBName:        getEnv("DB_NAME", "ottdb"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		Port:          getEnv("PORT", "8080"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@ott.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
