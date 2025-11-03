package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // ignore error; fallback to env

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBUrl:     getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/meditrack?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "change_me_please"),
	}
	log.Printf("Loaded config: port=%s\n", cfg.Port)
	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
