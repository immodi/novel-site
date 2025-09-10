package config

import (
	"log"
	"os"
)

var (
	DBPath       = mustGetEnv("DB_PATH")
	JWTSecret    = mustGetEnv("JWT_SECRET")
	Port         = mustGetEnv("PORT")
	IsProduction = mustGetEnv("IS_PRODUCTION") == "true"
)

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return value
}

