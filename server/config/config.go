package config

import (
	"os"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	HTTPSMode     string
	Domain        string
	Port          string
	AdminUsername string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgresql://admin:password@db:5432/consumable_db?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_to_a_random_secret_key_32_chars"),
		HTTPSMode:     getEnv("HTTPS_MODE", "selfsigned"),
		Domain:        getEnv("DOMAIN", "localhost"),
		Port:          getEnv("PORT", "8443"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}
}

func (c *Config) IsHTTPS() bool {
	return c.HTTPSMode != "http"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
