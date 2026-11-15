// Package config membaca konfigurasi aplikasi dari environment variable.
package config

import (
	"errors"
	"os"
	"strconv"
)

// Config berisi seluruh konfigurasi runtime aplikasi.
type Config struct {
	DatabaseURL   string
	Port          string
	SessionSecret string
	AdminPassword string
	SeedOnBoot    bool
}

// Load membaca konfigurasi dari environment variable.
func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:   env("DATABASE_URL", "postgres://dina:dina@localhost:5432/dinacatering?sslmode=disable"),
		Port:          env("PORT", "8080"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		AdminPassword: env("ADMIN_PASSWORD", "admin123"),
		SeedOnBoot:    envBool("SEED_ON_BOOT", true),
	}
	if cfg.SessionSecret == "" {
		return cfg, errors.New("SESSION_SECRET wajib diisi (lihat .env.example)")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}
