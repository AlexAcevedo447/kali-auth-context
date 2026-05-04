package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	AppName string
	AppEnv  string
	AppPort string

	// Database
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	// JWT
	JWTIssuer                string
	JWTAudience              string
	JWTSecret                string
	JWTAccessTokenTTLMinutes int

	// Keycloak
	KeycloakURL   string
	KeycloakRealm string

	// Seed
	SeedMasterEnabled              bool
	SeedMasterTenantID             string
	SeedMasterTenantName           string
	SeedMasterUserID               string
	SeedMasterIdentificationNumber string
	SeedMasterUsername             string
	SeedMasterEmail                string
	SeedMasterPassword             string
	SeedMasterPasswordFile         string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSLMODE"),

		JWTIssuer:                os.Getenv("JWT_ISSUER"),
		JWTAudience:              os.Getenv("JWT_AUDIENCE"),
		JWTSecret:                os.Getenv("JWT_SECRET"),
		JWTAccessTokenTTLMinutes: getEnvAsIntOrDefault("JWT_ACCESS_TOKEN_TTL_MINUTES", 60),

		KeycloakURL:   os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm: os.Getenv("KEYCLOAK_REALM"),

		SeedMasterEnabled:              getEnvAsBoolOrDefault("SEED_MASTER_ENABLED", false),
		SeedMasterTenantID:             getEnvOrDefault("SEED_MASTER_TENANT_ID", "master-tenant"),
		SeedMasterTenantName:           getEnvOrDefault("SEED_MASTER_TENANT_NAME", "Master Tenant"),
		SeedMasterUserID:               getEnvOrDefault("SEED_MASTER_USER_ID", "master-user"),
		SeedMasterIdentificationNumber: getEnvOrDefault("SEED_MASTER_IDENTIFICATION_NUMBER", "1000000000"),
		SeedMasterUsername:             getEnvOrDefault("SEED_MASTER_USERNAME", "master"),
		SeedMasterEmail:                getEnvOrDefault("SEED_MASTER_EMAIL", "master@local.dev"),
		SeedMasterPassword:             os.Getenv("SEED_MASTER_PASSWORD"),
		SeedMasterPasswordFile:         os.Getenv("SEED_MASTER_PASSWORD_FILE"),
	}
}

func (c *Config) ResolveSeedMasterPassword() (string, error) {
	if !c.SeedMasterEnabled {
		return "", nil
	}

	if c.SeedMasterPasswordFile != "" {
		bytes, err := os.ReadFile(c.SeedMasterPasswordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read SEED_MASTER_PASSWORD_FILE: %w", err)
		}

		secret := strings.TrimSpace(string(bytes))
		if secret == "" {
			return "", fmt.Errorf("SEED_MASTER_PASSWORD_FILE is empty")
		}

		return secret, nil
	}

	if strings.TrimSpace(c.SeedMasterPassword) == "" {
		return "", fmt.Errorf("SEED_MASTER_PASSWORD or SEED_MASTER_PASSWORD_FILE is required when SEED_MASTER_ENABLED=true")
	}

	return c.SeedMasterPassword, nil
}

func (c *Config) getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func (c *Config) getEnvAsInt(key string) int {
	valueStr := os.Getenv(key)

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid integer", key)
	}

	return value
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid integer", key)
	}

	return value
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid boolean", key)
	}

	return parsed
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}