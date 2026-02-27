// Package config demonstrates how to read environment variables, JSON, and YAML
// in a Go application — a complete, production-ready configuration pattern.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ============================================================
// PART 1: ENVIRONMENT VARIABLES
// ============================================================
// os.Getenv, os.LookupEnv, os.Setenv — no external deps needed.

// EnvExample shows all the ways to read environment variables.
func EnvExample() {
	// ── Basic read ──────────────────────────────────────────
	// Returns "" if not set — you can't distinguish "not set" from ""
	port := os.Getenv("APP_PORT")
	fmt.Println("PORT:", port)

	// ── Safe read with existence check ──────────────────────
	// LookupEnv returns (value, exists bool)
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		fmt.Println("DATABASE_URL is not set — using default")
		dbURL = "postgres://localhost:5432/mydb"
	}
	fmt.Println("DB:", dbURL)

	// ── Parse int from env ───────────────────────────────────
	portStr := os.Getenv("APP_PORT")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		portInt = 8080 // fallback default
	}
	fmt.Println("Port as int:", portInt)

	// ── Parse bool from env ──────────────────────────────────
	debugStr := os.Getenv("DEBUG")
	debug, _ := strconv.ParseBool(debugStr) // "true", "1", "yes" → true
	fmt.Println("Debug mode:", debug)

	// ── Set an env var (within this process only) ────────────
	os.Setenv("MY_VAR", "hello")
	fmt.Println("MY_VAR:", os.Getenv("MY_VAR"))

	// ── List all environment variables ───────────────────────
	// envVars := os.Environ() // returns []string of "KEY=VALUE"
}

// GetEnvOrDefault returns the env variable or a fallback default.
// This is a very common helper pattern in Go projects.
func GetEnvOrDefault(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}

// GetEnvIntOrDefault returns an int env variable or default.
func GetEnvIntOrDefault(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

// ============================================================
// PART 2: .ENV FILE (godotenv)
// ============================================================
// In production, env vars come from the system (k8s secrets, etc.)
// In development, use a .env file with the godotenv library.
//
// Install: go get github.com/joho/godotenv
//
// ─── .env file (add to .gitignore!) ──────────────────────────
//   APP_PORT=8080
//   APP_ENV=development
//   DATABASE_URL=postgres://localhost:5432/mydev
//   JWT_SECRET=my-super-secret-key
//   DEBUG=true
// ─────────────────────────────────────────────────────────────
//
// Usage in main.go:
//
//   import "github.com/joho/godotenv"
//
//   func main() {
//       // Load .env file — only in development, not in production
//       if err := godotenv.Load(); err != nil {
//           log.Println("No .env file found, using system env vars")
//       }
//       port := os.Getenv("APP_PORT")
//       // ... rest of app
//   }
//
// Load a specific file:
//   godotenv.Load(".env.local")
//   godotenv.Load(".env", ".env.local")  // later file overrides earlier
//
// Overload (override existing env vars):
//   godotenv.Overload(".env.local")

// ============================================================
// PART 3: JSON CONFIG FILE
// ============================================================

// AppConfig is the top-level config struct — maps 1:1 to your JSON/YAML.
// json tags tell the JSON decoder which key maps to which field.
// `json:"key,omitempty"` — omit the field if it's the zero value when marshaling.
type AppConfig struct {
	App      AppSettings      `json:"app"`
	Database DatabaseSettings `json:"database"`
	Cache    CacheSettings    `json:"cache"`
	Features FeatureFlags     `json:"features"`
}

// AppSettings holds general app configuration.
type AppSettings struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Port        int    `json:"port"`
	Debug       bool   `json:"debug"`
	LogLevel    string `json:"log_level"`
}

// DatabaseSettings holds database configuration.
type DatabaseSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"` // in real apps, read from env not config file
	SSLMode  string `json:"ssl_mode"`
	MaxConns int    `json:"max_connections"`
}

// DSN builds a Postgres DSN from settings.
func (d DatabaseSettings) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// CacheSettings holds Redis configuration.
type CacheSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	TTLSecs  int    `json:"ttl_seconds"`
}

// FeatureFlags enables/disables features without redeployment.
type FeatureFlags struct {
	EnableNewUI      bool `json:"enable_new_ui"`
	EnableBetaSearch bool `json:"enable_beta_search"`
	RateLimit        bool `json:"rate_limit"`
}

// LoadJSONConfig reads and parses a JSON config file.
// Usage: cfg, err := LoadJSONConfig("configs/config.json")
func LoadJSONConfig(path string) (*AppConfig, error) {
	// Read the entire file into memory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg AppConfig
	// json.Unmarshal parses JSON bytes into the struct
	// Field names are matched by json tags (case-insensitive fallback)
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	return &cfg, nil
}

// ─── configs/config.json ─────────────────────────────────────
// {
//   "app": {
//     "name": "my-go-service",
//     "environment": "development",
//     "port": 8080,
//     "debug": true,
//     "log_level": "debug"
//   },
//   "database": {
//     "host": "localhost",
//     "port": 5432,
//     "name": "mydb",
//     "user": "postgres",
//     "password": "secret",
//     "ssl_mode": "disable",
//     "max_connections": 25
//   },
//   "cache": {
//     "host": "localhost",
//     "port": 6379,
//     "password": "",
//     "db": 0,
//     "ttl_seconds": 300
//   },
//   "features": {
//     "enable_new_ui": true,
//     "enable_beta_search": false,
//     "rate_limit": true
//   }
// }
// ─────────────────────────────────────────────────────────────

// SaveJSONConfig serializes config to a JSON file (useful for defaults).
func SaveJSONConfig(path string, cfg *AppConfig) error {
	// MarshalIndent = pretty-printed JSON with 2-space indent
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing config: %w", err)
	}
	// 0644 = owner read/write, group+others read-only
	return os.WriteFile(path, data, 0644)
}

// ============================================================
// PART 4: YAML CONFIG FILE
// ============================================================
// Go has no built-in YAML parser. Use gopkg.in/yaml.v3.
// Install: go get gopkg.in/yaml.v3
//
// yaml tags work exactly like json tags.
// The struct below is the same AppConfig but with yaml tags added.

// AppConfigYAML is the same config with YAML tags.
// In real projects, you'd add BOTH json and yaml tags to a single struct:
//
//	Port int `json:"port" yaml:"port"`
type AppConfigYAML struct {
	App      AppSettingsYAML      `yaml:"app"`
	Database DatabaseSettingsYAML `yaml:"database"`
	Features FeatureFlagsYAML     `yaml:"features"`
}

type AppSettingsYAML struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Port        int    `yaml:"port"`
	Debug       bool   `yaml:"debug"`
	LogLevel    string `yaml:"log_level"`
}

type DatabaseSettingsYAML struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

type FeatureFlagsYAML struct {
	EnableNewUI      bool `yaml:"enable_new_ui"`
	EnableBetaSearch bool `yaml:"enable_beta_search"`
}

// ─── configs/config.yaml ─────────────────────────────────────
// app:
//   name: my-go-service
//   environment: development
//   port: 8080
//   debug: true
//   log_level: debug
//
// database:
//   host: localhost
//   port: 5432
//   name: mydb
//   user: postgres
//   password: secret       # use env var in production!
//   ssl_mode: disable
//
// features:
//   enable_new_ui: true
//   enable_beta_search: false
// ─────────────────────────────────────────────────────────────
//
// LoadYAMLConfig reads and parses a YAML config file.
// Uncomment and use when gopkg.in/yaml.v3 is added to go.mod:
//
// import "gopkg.in/yaml.v3"
//
// func LoadYAMLConfig(path string) (*AppConfigYAML, error) {
//     data, err := os.ReadFile(path)
//     if err != nil {
//         return nil, fmt.Errorf("reading yaml config %q: %w", path, err)
//     }
//     var cfg AppConfigYAML
//     if err := yaml.Unmarshal(data, &cfg); err != nil {
//         return nil, fmt.Errorf("parsing yaml config %q: %w", path, err)
//     }
//     return &cfg, nil
// }

// ============================================================
// PART 5: COMPLETE PRODUCTION-READY CONFIG PATTERN
// ============================================================
// Best practice: combine env vars + config file.
// Env vars OVERRIDE config file values (12-factor app principle).

// Config is the single source of truth for the whole application.
// Pass this struct (or a pointer to it) through your dependency tree.
type Config struct {
	Env      string // "development", "staging", "production"
	Port     int
	LogLevel string
	DB       DBConfig
	Cache    RedisConfig
}

type DBConfig struct {
	URL      string
	MaxConns int
}

type RedisConfig struct {
	URL string
}

// Load builds Config from env vars (primary) + JSON file (fallback defaults).
// This is the canonical pattern:
//  1. Load defaults from config file
//  2. Override with environment variables
//  3. Validate required fields
func Load(configPath string) (*Config, error) {
	// Step 1: load from JSON file as defaults
	var fileConf AppConfig
	if data, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(data, &fileConf) // ignore error — env vars will cover it
	}

	// Step 2: build config, env vars win over file values
	cfg := &Config{
		Env:      GetEnvOrDefault("APP_ENV", fileConf.App.Environment),
		Port:     GetEnvIntOrDefault("APP_PORT", fileConf.App.Port),
		LogLevel: GetEnvOrDefault("LOG_LEVEL", fileConf.App.LogLevel),
		DB: DBConfig{
			// DATABASE_URL env var completely overrides file settings
			URL:      GetEnvOrDefault("DATABASE_URL", fileConf.Database.DSN()),
			MaxConns: GetEnvIntOrDefault("DB_MAX_CONNS", fileConf.Database.MaxConns),
		},
		Cache: RedisConfig{
			URL: GetEnvOrDefault("REDIS_URL",
				fmt.Sprintf("redis://%s:%d", fileConf.Cache.Host, fileConf.Cache.Port)),
		},
	}

	// Step 3: apply defaults for required fields
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.Env == "" {
		cfg.Env = "development"
	}

	// Step 4: validate
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DB.URL == "" {
		return errors.New("DATABASE_URL is required")
	}
	return nil
}

// IsDevelopment returns true in development mode.
func (c *Config) IsDevelopment() bool { return c.Env == "development" }

// IsProduction returns true in production mode.
func (c *Config) IsProduction() bool { return c.Env == "production" }

// ============================================================
// PART 6: USING CONFIG IN main.go  (pattern)
// ============================================================
//
//   func main() {
//       // Load .env in development only
//       if os.Getenv("APP_ENV") != "production" {
//           godotenv.Load() // silently ignore if .env doesn't exist
//       }
//
//       cfg, err := config.Load("configs/config.json")
//       if err != nil {
//           log.Fatalf("failed to load config: %v", err)
//       }
//
//       slog.SetLogLoggerLevel(slog.LevelDebug)
//       if cfg.IsProduction() {
//           slog.SetLogLoggerLevel(slog.LevelInfo)
//       }
//
//       db, err := sql.Open("pgx", cfg.DB.URL)
//       if err != nil {
//           log.Fatalf("failed to connect to DB: %v", err)
//       }
//
//       server := NewServer(cfg, db)
//       log.Printf("starting server on :%d (env=%s)", cfg.Port, cfg.Env)
//       log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), server))
//   }
