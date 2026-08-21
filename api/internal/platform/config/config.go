package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envFileStateMu      sync.Mutex
	loadedEnvFileValues = make(map[string]string)
)

type Config struct {
	Version         string
	Environment     string
	LogLevel        string
	APIPort         string
	FrontendOrigin  string
	CookieSecure    bool
	Database        DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (*Config, error) {
	environment := getEnv("ENVIRONMENT", "development")
	if err := loadEnvFiles(environment); err != nil {
		return nil, fmt.Errorf("load env files error: %w", err)
	}

	missing := make([]string, 0, 7)
	dbHost := mustEnv("DB_HOST", &missing)
	dbPort := mustEnv("DB_PORT", &missing)
	dbUser := mustEnv("DB_USER", &missing)
	dbPassword := mustEnv("DB_PASSWORD", &missing)
	dbName := mustEnv("DB_NAME", &missing)
	dbSSLMode := mustEnv("DB_SSL_MODE", &missing)
	frontendOrigin := mustEnv("FRONTEND_ORIGIN", &missing)
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	apiPort := getEnv("API_PORT", "8180")
	if _, err := strconv.Atoi(apiPort); err != nil {
		return nil, fmt.Errorf("invalid value for API_PORT: %w", err)
	}

	return &Config{
		Version:        getEnv("VERSION", "dev"),
		Environment:    environment,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		APIPort:        apiPort,
		FrontendOrigin: frontendOrigin,
		CookieSecure:   environment == "production" || environment == "prod",
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if raw, ok := os.LookupEnv(key); ok && raw != "" {
		return strings.TrimSpace(raw)
	}
	return defaultValue
}

func mustEnv(key string, missing *[]string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		*missing = append(*missing, key)
	}
	return value
}

func envSearchPaths() []string {
	paths := make([]string, 0, 8)
	seen := make(map[string]struct{})
	add := func(path string) {
		cleaned := filepath.Clean(path)
		if _, ok := seen[cleaned]; ok {
			return
		}
		seen[cleaned] = struct{}{}
		paths = append(paths, cleaned)
	}

	if wd, err := os.Getwd(); err == nil {
		for dir := wd; ; dir = filepath.Dir(dir) {
			add(dir)
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	add(".")
	add("..")
	return paths
}

func loadEnvFiles(environment string) error {
	envName := strings.ToLower(strings.TrimSpace(environment))
	envSpecificFile := fmt.Sprintf(".env.%s", envName)
	searchPaths := envSearchPaths()

	merged := make(map[string]string)
	baseValues := make(map[string]string)
	environmentValues := make(map[string]string)

	for _, dir := range searchPaths {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			values, err := godotenv.Read(path)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			for k, v := range values {
				baseValues[k] = v
				merged[k] = v
			}
			break
		}
	}

	for _, dir := range searchPaths {
		path := filepath.Join(dir, envSpecificFile)
		if _, err := os.Stat(path); err == nil {
			values, err := godotenv.Read(path)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			for k, v := range values {
				environmentValues[k] = v
				merged[k] = v
			}
			break
		}
	}

	envFileStateMu.Lock()
	defer envFileStateMu.Unlock()
	for k, v := range merged {
		current, exists := os.LookupEnv(k)
		previousLoaded, managed := loadedEnvFileValues[k]
		if exists && !managed {
			baseValue, hasBaseValue := baseValues[k]
			_, hasEnvironmentValue := environmentValues[k]
			if !hasEnvironmentValue || !hasBaseValue || current != baseValue {
				continue
			}
		}
		if exists && managed && current != previousLoaded {
			delete(loadedEnvFileValues, k)
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("setenv error: %w", err)
		}
		loadedEnvFileValues[k] = v
	}

	return nil
}
