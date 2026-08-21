package config

import (
	"os"
	"testing"
)

func isolateEnvFiles(t testing.TB) {
	t.Helper()
	dir := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to capture cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(original); chdirErr != nil {
			t.Errorf("failed to restore cwd: %v", chdirErr)
		}
	})
}

func clearConfigEnv(t testing.TB) {
	t.Helper()
	envFileStateMu.Lock()
	loadedEnvFileValues = make(map[string]string)
	envFileStateMu.Unlock()
	for _, key := range []string{
		"ENVIRONMENT", "VERSION", "LOG_LEVEL", "API_PORT", "FRONTEND_ORIGIN",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
	} {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset %s: %v", key, err)
		}
	}
}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "flower")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "flower_test")
	t.Setenv("DB_SSL_MODE", "disable")
	t.Setenv("FRONTEND_ORIGIN", "http://localhost:4273")
}

func TestLoadRejectsMissingDatabaseSettings(t *testing.T) {
	isolateEnvFiles(t)
	clearConfigEnv(t)
	t.Setenv("ENVIRONMENT", "test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing required environment variables to fail Load")
	}
	if got := err.Error(); got == "" {
		t.Fatal("expected a descriptive error, got empty string")
	}
}

func TestLoadReadsRequiredDatabaseSettings(t *testing.T) {
	isolateEnvFiles(t)
	clearConfigEnv(t)
	setRequiredEnv(t)
	t.Setenv("ENVIRONMENT", "test")
	t.Setenv("VERSION", "1.2.3")
	t.Setenv("API_PORT", "8099")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Version != "1.2.3" {
		t.Fatalf("version: got %q, want %q", cfg.Version, "1.2.3")
	}
	if cfg.Environment != "test" {
		t.Fatalf("environment: got %q, want %q", cfg.Environment, "test")
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("log level: got %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.APIPort != "8099" {
		t.Fatalf("api port: got %q, want %q", cfg.APIPort, "8099")
	}
	if cfg.Database.Host != "localhost" {
		t.Fatalf("db host: got %q, want %q", cfg.Database.Host, "localhost")
	}
	if cfg.Database.Port != "5432" {
		t.Fatalf("db port: got %q, want %q", cfg.Database.Port, "5432")
	}
	if cfg.Database.User != "flower" {
		t.Fatalf("db user: got %q, want %q", cfg.Database.User, "flower")
	}
	if cfg.Database.Password != "secret" {
		t.Fatalf("db password: got %q, want %q", cfg.Database.Password, "secret")
	}
	if cfg.Database.Name != "flower_test" {
		t.Fatalf("db name: got %q, want %q", cfg.Database.Name, "flower_test")
	}
	if cfg.Database.SSLMode != "disable" {
		t.Fatalf("db ssl mode: got %q, want %q", cfg.Database.SSLMode, "disable")
	}
	if cfg.FrontendOrigin != "http://localhost:4273" {
		t.Fatalf("frontend origin: got %q", cfg.FrontendOrigin)
	}
	if cfg.CookieSecure {
		t.Fatal("cookie secure should be false in test")
	}
}

func TestLoadDefaultsAPIPortToFlowerHostPort(t *testing.T) {
	isolateEnvFiles(t)
	clearConfigEnv(t)
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.APIPort != "8180" {
		t.Fatalf("api port: got %q, want %q", cfg.APIPort, "8180")
	}
}

func TestLoadRejectsInvalidAPIPort(t *testing.T) {
	isolateEnvFiles(t)
	clearConfigEnv(t)
	setRequiredEnv(t)
	t.Setenv("API_PORT", "not-a-port")

	_, err := Load()
	if err == nil {
		t.Fatal("expected invalid API_PORT to fail Load")
	}
}
