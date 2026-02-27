package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue string
		want         string
	}{
		{"env set", "TEST_KEY", "myvalue", true, "default", "myvalue"},
		{"env not set", "TEST_KEY_MISSING", "", false, "default", "default"},
		{"env empty uses default", "TEST_KEY_EMPTY", "", true, "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}
			got := GetEnvOrDefault(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnvOrDefault(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue int
		want         int
	}{
		{"valid int", "TEST_PORT", "9090", true, 8080, 9090},
		{"invalid int falls back", "TEST_PORT_BAD", "notanint", true, 8080, 8080},
		{"not set falls back", "TEST_PORT_MISSING", "", false, 8080, 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}
			got := GetEnvIntOrDefault(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnvIntOrDefault(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}

func TestLoadJSONConfig(t *testing.T) {
	cfg, err := LoadJSONConfig("config.json")
	if err != nil {
		t.Fatalf("LoadJSONConfig() error = %v", err)
	}

	if cfg.App.Name != "go-interview-prep" {
		t.Errorf("App.Name = %q, want %q", cfg.App.Name, "go-interview-prep")
	}
	if cfg.App.Port != 8080 {
		t.Errorf("App.Port = %d, want %d", cfg.App.Port, 8080)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "localhost")
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Database.Port = %d, want %d", cfg.Database.Port, 5432)
	}
	if !cfg.Features.EnableNewUI {
		t.Error("Features.EnableNewUI should be true")
	}
}

func TestLoadJSONConfigMissing(t *testing.T) {
	_, err := LoadJSONConfig("nonexistent.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestDatabaseDSN(t *testing.T) {
	db := DatabaseSettings{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		Name:     "mydb",
		SSLMode:  "disable",
	}
	dsn := db.DSN()
	if dsn == "" {
		t.Error("DSN() should not be empty")
	}
	// Check it contains key parts
	for _, part := range []string{"localhost", "postgres", "mydb", "disable"} {
		found := false
		for i := 0; i+len(part) <= len(dsn); i++ {
			if dsn[i:i+len(part)] == part {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DSN() missing %q in result: %q", part, dsn)
		}
	}
}

func TestConfigEnvOverride(t *testing.T) {
	// Set env vars to override config file
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("DATABASE_URL", "postgres://prod-host:5432/proddb")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg, err := Load("config.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != "production" {
		t.Errorf("Env = %q, want %q", cfg.Env, "production")
	}
	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 9000)
	}
	if cfg.DB.URL != "postgres://prod-host:5432/proddb" {
		t.Errorf("DB.URL = %q, want %q", cfg.DB.URL, "postgres://prod-host:5432/proddb")
	}
	if !cfg.IsProduction() {
		t.Error("IsProduction() should be true")
	}
	if cfg.IsDevelopment() {
		t.Error("IsDevelopment() should be false in production")
	}
}
