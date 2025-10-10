package config

import (
	"log"
	"os"
)

var (
	DBPath             = mustGetEnv("DB_PATH")
	JWTSecret          = mustGetEnv("JWT_SECRET")
	Port               = mustGetEnv("PORT")
	SiteURL            = mustGetEnv("SITE_URL")
	AdminSiteURL       = mustGetEnv("ADMIN_SITE_URL")
	IsProduction       = mustGetEnv("IS_PRODUCTION") == "true"
	GoogleClientID     = mustGetEnv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = mustGetEnv("GOOGLE_CLIENT_SECRET")
	GoogleFormURL      = mustGetEnv("GOOGLE_FORM_URL")
)

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return value
}
