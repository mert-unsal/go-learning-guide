# 20 — Practical Go Toolchain

> Dependency management, building, Docker, and configuration — the operational side of shipping Go to production.

---

## Table of Contents

1. [Go Modules & Dependencies](#1-go-modules--dependencies)
2. [Build, Run & Deploy](#2-build-run--deploy)
3. [Dockerizing Go Applications](#3-dockerizing-go-applications)
4. [Configuration & Environment](#4-configuration--environment)
5. [Quick Reference Cheatsheet](#5-quick-reference-cheatsheet)

---

## 1. Go Modules & Dependencies

A Go module is a collection of packages versioned together, defined by `go.mod` at the root.

### Initialise a Module

```bash
mkdir myapp && cd myapp
go mod init github.com/yourname/myapp
```

This creates a minimal `go.mod`:

```
module github.com/yourname/myapp
go 1.25.7
```

### Adding Dependencies

```bash
go get github.com/some/package           # latest version
go get github.com/some/package@v1.2.3    # exact version
go get github.com/some/package@latest    # explicitly latest
```

After `go get`, both `go.mod` (require line) and `go.sum` (hashes) are updated.

### Reading go.mod

```
module github.com/yourname/myapp

go 1.25.7

require (
    github.com/gin-gonic/gin v1.9.1          // direct dependency
    gopkg.in/yaml.v3 v3.0.1
)

require (
    golang.org/x/net v0.17.0 // indirect      // pulled in by a direct dep
)

replace github.com/some/pkg => ../local-fork   // local override for dev
```

| Directive   | Purpose                                          |
|-------------|--------------------------------------------------|
| `module`    | Declares the module path (import path root)      |
| `go`        | Minimum Go version                               |
| `require`   | Lists dependencies — direct and `// indirect`    |
| `replace`   | Swap a dependency for a local path or fork       |
| `exclude`   | Prevent a specific version from being used       |
| `retract`   | Mark your own module versions as broken           |

### go.sum — What It Is

- Auto-generated SHA-256 hashes for every dependency version.
- **Never edit manually.** Never `.gitignore` it.
- **Always commit it** — ensures reproducible builds. CI and teammates get identical dependency trees.

### Module Management Commands

| Command            | What It Does                                              |
|--------------------|-----------------------------------------------------------|
| `go mod tidy`      | Remove unused deps, add missing ones — run after every change |
| `go mod vendor`    | Copy all deps into `./vendor/` (for air-gapped builds)    |
| `go mod verify`    | Check deps haven't been tampered with (hash integrity)    |
| `go mod download`  | Pre-download all deps to local cache (`$GOPATH/pkg/mod`)  |
| `go mod graph`     | Print the full dependency graph (pipe to `grep` to filter)|

### Upgrade & Downgrade Strategies

```bash
go get github.com/pkg@v1.3.0      # pin exact version
go get github.com/pkg@v1           # latest v1.x.x
go get -u ./...                    # upgrade ALL deps (careful — may break things)
go get -u=patch ./...              # patch bumps only (safest bulk upgrade)
```

To remove a dependency: delete all imports from your code, then `go mod tidy`.

### Workspace Mode (go.work) — Monorepo Development

```bash
go work init ./serviceA ./serviceB
go work use ./shared-lib
```

Creates `go.work` — lets you develop multiple modules that depend on each other without `replace` directives. **Do not commit `go.work`** — add it to `.gitignore`.

### Private Modules

```bash
export GOPRIVATE=github.com/yourorg/*      # skip proxy + sum DB
export GONOSUMCHECK=github.com/yourorg/*   # skip checksum verification
export GONOSUMDB=github.com/yourorg/*      # skip sum database lookup
```

> **Key insight:** `GOPRIVATE` is a superset — setting it also implies `GONOSUMCHECK` and `GONOSUMDB` for matching patterns.

### Common Dependencies

| Package                          | Purpose               | Install                                    |
|----------------------------------|-----------------------|--------------------------------------------|
| `github.com/joho/godotenv`       | `.env` file loading   | `go get github.com/joho/godotenv`          |
| `gopkg.in/yaml.v3`              | YAML parsing          | `go get gopkg.in/yaml.v3`                  |
| `github.com/gin-gonic/gin`      | HTTP web framework    | `go get github.com/gin-gonic/gin`          |
| `go.uber.org/zap`               | Structured logging    | `go get go.uber.org/zap`                   |
| `github.com/stretchr/testify`   | Test assertions/mocks | `go get github.com/stretchr/testify`       |
| `gorm.io/gorm`                  | Database ORM          | `go get gorm.io/gorm`                      |

---

## 2. Build, Run & Deploy

### go run vs go build vs go install

| Command        | What It Does                                          | Use When                        |
|----------------|-------------------------------------------------------|---------------------------------|
| `go run .`     | Compile + execute in one step (temp binary, deleted)  | Development iteration           |
| `go build .`   | Compile to a binary in current directory              | Building for deployment         |
| `go install .` | Compile + place binary in `$GOPATH/bin`               | Installing CLI tools globally   |

```bash
go run main.go                      # single file
go run .                            # all files in current dir
go run ./cmd/myapp/                  # specific entry point
go build -o myapp ./cmd/myapp/      # named output binary
go install github.com/tool@latest   # install a CLI tool
```

### Build Flags

| Flag              | Purpose                                   | Example                           |
|-------------------|-------------------------------------------|-----------------------------------|
| `-o <name>`       | Output binary name                        | `-o bin/myapp`                    |
| `-v`              | Verbose — show packages being compiled    | `-v`                              |
| `-race`           | Enable race detector (use in tests!)      | `-race`                           |
| `-gcflags`        | Compiler flags (escape analysis, debug)   | `-gcflags="-m -m"` / `"-N -l"`   |
| `-ldflags`        | Linker flags (version injection, strip)   | `-ldflags "-s -w"`               |
| `-tags`           | Build tags for conditional compilation    | `-tags integration`               |

### Version Injection with ldflags

Declare variables in your code, then set them at build time — no config file needed:

```go
// main.go
var (
    version   = "dev"
    buildTime = "unknown"
    gitCommit = "unknown"
)
```

```bash
go build \
  -ldflags "-X main.version=1.0.0 \
            -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
            -X main.gitCommit=$(git rev-parse --short HEAD)" \
  -o myapp ./cmd/myapp/
```

### Cross-Compilation

Go cross-compiles natively — no extra toolchain needed. Set `GOOS` and `GOARCH`:

```bash
GOOS=linux   GOARCH=amd64 go build -o myapp-linux .
GOOS=windows GOARCH=amd64 go build -o myapp.exe .
GOOS=darwin  GOARCH=arm64 go build -o myapp-mac .
```

| GOOS      | GOARCH  | Target                        |
|-----------|---------|-------------------------------|
| `linux`   | `amd64` | Linux x86-64 (servers, CI)    |
| `linux`   | `arm64` | Linux ARM (AWS Graviton, RPi) |
| `darwin`  | `arm64` | macOS Apple Silicon           |
| `darwin`  | `amd64` | macOS Intel                   |
| `windows` | `amd64` | Windows x86-64                |

Full list: `go tool dist list`

### Strip Debug Info (Smaller Binary)

```bash
go build -ldflags "-s -w" -o myapp .
# -s: omit symbol table
# -w: omit DWARF debug info
# Result: ~30% smaller binary
```

> **Tradeoff:** Stripped binaries can't be debugged with `dlv`. Use this for production images, not development. See [Chapter 15](./15_debugging_profiling.md) for debugging workflows.

### Project Structure

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go           # entry point (package main)
├── internal/                 # private packages (compiler-enforced)
│   ├── config/
│   ├── handler/
│   └── service/
├── pkg/                      # public reusable packages (use sparingly)
├── configs/
│   ├── config.yaml
│   └── config.local.yaml
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── go.sum
```

### Makefile Template

```makefile
.PHONY: build run test clean lint

build:
	go build -ldflags "-s -w" -o bin/myapp ./cmd/myapp/

run:
	go run ./cmd/myapp/

test:
	go test -race -cover ./...

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...
```

### Systemd Service File

For Linux bare-metal or VM deployments:

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Go App
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/myapp
EnvironmentFile=/opt/myapp/.env
ExecStart=/opt/myapp/myapp
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable myapp     # start on boot
systemctl start  myapp     # start now
systemctl status myapp     # check status
journalctl -u myapp -f     # tail logs
```

---

## 3. Dockerizing Go Applications

### Simple Dockerfile (Learning)

Good for understanding — **not for production** (image is 300+ MB):

```dockerfile
FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/myapp/

EXPOSE 8080
CMD ["./myapp"]
```

### Production Dockerfile — Multi-Stage Build

This is what you ship. Every line exists for a reason:

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git              # needed if deps use private repos

COPY go.mod go.sum ./
RUN go mod download                      # cached layer — only re-runs if go.mod changes

COPY . .

# CGO_ENABLED=0: static binary, no libc dependency → runs on scratch
# -ldflags "-s -w": strip debug info → smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o myapp ./cmd/myapp/

# ── Stage 2: Final image ──────────────────────────────────
FROM scratch
# Alternatives:
#   FROM gcr.io/distroless/static:nonroot  ← Google's minimal image (has CA certs)
#   FROM alpine:3.19                       ← if you need a shell for debugging

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/configs/ /configs/

USER 65534:65534                          # non-root (nobody user)

EXPOSE 8080

ENTRYPOINT ["/myapp"]
```

| Base Image Choice     | Size    | Shell | CA Certs | Use Case                   |
|-----------------------|---------|-------|----------|----------------------------|
| `scratch`             | ~0 MB   | ❌    | ❌       | Minimal — copy certs yourself |
| `distroless/static`   | ~2 MB   | ❌    | ✅       | Production default           |
| `alpine`              | ~7 MB   | ✅    | ✅       | Need shell for debugging     |

> **Why `ENTRYPOINT` not `CMD`?** `ENTRYPOINT` makes the container behave like an executable. `CMD` provides default args that can be overridden. Use `ENTRYPOINT` for the binary, `CMD` for default flags.

### .dockerignore

Place next to your Dockerfile — keeps build context small and secure:

```
.git
.gitignore
*.md
bin/
vendor/
.env
.env.local
*_test.go
```

### Docker Commands

| Task               | Command                                                    |
|--------------------|------------------------------------------------------------|
| Build              | `docker build -t myapp:latest .`                           |
| Build (tagged)     | `docker build -t myapp:1.0.0 .`                           |
| Build (no cache)   | `docker build --no-cache -t myapp:latest .`                |
| Run                | `docker run -p 8080:8080 myapp:latest`                     |
| Run (background)   | `docker run -d -p 8080:8080 --name myapp myapp:latest`    |
| Run (with env)     | `docker run -e APP_ENV=prod --env-file .env myapp:latest`  |
| Run (with volume)  | `docker run -v $(pwd)/configs:/configs myapp:latest`       |
| Logs               | `docker logs -f myapp`                                     |
| Shell into         | `docker exec -it myapp sh`                                 |
| Stop               | `docker stop myapp`                                        |
| Remove             | `docker rm myapp`                                          |
| List running       | `docker ps`                                                |
| List all           | `docker ps -a`                                             |

### Docker Compose — Multi-Service Local Dev

```yaml
version: "3.9"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - DB_URL=postgres://postgres:password@db:5432/mydb
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - ./configs:/configs

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

### Docker Compose Commands

| Task                | Command                            |
|---------------------|------------------------------------|
| Start all           | `docker compose up`                |
| Start (background)  | `docker compose up -d`             |
| Start (rebuild)     | `docker compose up --build`        |
| Stop + remove       | `docker compose down`              |
| Stop + remove vols  | `docker compose down -v`           |
| Follow logs         | `docker compose logs -f app`       |
| Service status      | `docker compose ps`                |
| Shell into service  | `docker compose exec app sh`       |

### Push to Registry

```bash
# Docker Hub
docker login
docker tag myapp:latest yourusername/myapp:latest
docker push yourusername/myapp:latest

# AWS ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker tag myapp:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/myapp:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/myapp:latest
```

### Hot Reload with Air (Development)

```dockerfile
# Dockerfile.dev
FROM golang:1.25-alpine
RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
CMD ["air"]
```

```bash
docker run -v $(pwd):/app -p 8080:8080 myapp-dev
```

> **Don't use air in production.** It watches for file changes and recompiles — development only.

---

## 4. Configuration & Environment

### os.Getenv vs os.LookupEnv

```go
// Getenv returns "" if not set — can't distinguish "not set" from "empty"
port := os.Getenv("APP_PORT")

// LookupEnv tells you whether the variable exists
dbURL, ok := os.LookupEnv("DATABASE_URL")
if !ok {
    dbURL = "postgres://localhost:5432/mydb"
}
```

> **Rule of thumb:** Use `os.LookupEnv` when you need to distinguish "not set" from "set to empty string." Use `os.Getenv` when you just need a value with a fallback.

### GetEnvOrDefault Helper — The Pattern Everyone Uses

```go
func GetEnvOrDefault(key, defaultValue string) string {
    if val, ok := os.LookupEnv(key); ok && val != "" {
        return val
    }
    return defaultValue
}

func GetEnvIntOrDefault(key string, defaultValue int) int {
    if val, ok := os.LookupEnv(key); ok {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultValue
}
```

### .env Files with godotenv (Dev Only)

```bash
# .env  — add to .gitignore!
APP_PORT=8080
APP_ENV=development
DATABASE_URL=postgres://localhost:5432/mydev
JWT_SECRET=my-super-secret-key
DEBUG=true
```

```go
import "github.com/joho/godotenv"

func main() {
    // Load .env — only in development, not production
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system env vars")
    }
    port := os.Getenv("APP_PORT")
}
```

> **In production, env vars come from the platform** — Kubernetes secrets, AWS Parameter Store, systemd `EnvironmentFile`. Never deploy `.env` files.

### JSON Config

Define structs with `json` tags — field names map to JSON keys:

```go
type AppConfig struct {
    App      AppSettings      `json:"app"`
    Database DatabaseSettings `json:"database"`
    Features FeatureFlags     `json:"features"`
}

type AppSettings struct {
    Name        string `json:"name"`
    Environment string `json:"environment"`
    Port        int    `json:"port"`
    Debug       bool   `json:"debug"`
    LogLevel    string `json:"log_level"`
}
```

```go
func LoadJSONConfig(path string) (*AppConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config %q: %w", path, err)
    }
    var cfg AppConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parsing config %q: %w", path, err)
    }
    return &cfg, nil
}
```

### YAML Config (gopkg.in/yaml.v3)

YAML tags work identically to JSON tags. In real projects, add both:

```go
type AppSettings struct {
    Name string `json:"name" yaml:"name"`
    Port int    `json:"port" yaml:"port"`
}
```

```yaml
# configs/config.yaml
app:
  name: my-go-service
  environment: development
  port: 8080
  debug: true

database:
  host: localhost
  port: 5432
  name: mydb
  user: postgres
  password: secret       # use env var in production!
  ssl_mode: disable
```

### Production Pattern: Env Vars Override Config File (12-Factor)

The canonical approach — config file provides defaults, env vars win:

```go
type Config struct {
    Env      string
    Port     int
    LogLevel string
    DB       DBConfig
}

func Load(configPath string) (*Config, error) {
    // Step 1: load file as defaults
    var fileConf AppConfig
    if data, err := os.ReadFile(configPath); err == nil {
        _ = json.Unmarshal(data, &fileConf)
    }

    // Step 2: env vars override file values
    cfg := &Config{
        Env:      GetEnvOrDefault("APP_ENV", fileConf.App.Environment),
        Port:     GetEnvIntOrDefault("APP_PORT", fileConf.App.Port),
        LogLevel: GetEnvOrDefault("LOG_LEVEL", fileConf.App.LogLevel),
    }

    // Step 3: apply defaults for zero values
    if cfg.Port == 0 {
        cfg.Port = 8080
    }

    // Step 4: validate required fields
    return cfg, cfg.validate()
}
```

### Validation Pattern

```go
func (c *Config) validate() error {
    if c.DB.URL == "" {
        return errors.New("DATABASE_URL is required")
    }
    return nil
}

func (c *Config) IsDevelopment() bool { return c.Env == "development" }
func (c *Config) IsProduction() bool  { return c.Env == "production" }
```

### Using Config in main.go

```go
func main() {
    // Load .env in development only
    if os.Getenv("APP_ENV") != "production" {
        godotenv.Load()
    }

    cfg, err := config.Load("configs/config.json")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    db, err := sql.Open("pgx", cfg.DB.URL)
    if err != nil {
        log.Fatalf("failed to connect to DB: %v", err)
    }

    server := NewServer(cfg, db)
    log.Printf("starting on :%d (env=%s)", cfg.Port, cfg.Env)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), server))
}
```

---

## 5. Quick Reference Cheatsheet

| Task                     | Command                                                         |
|--------------------------|-----------------------------------------------------------------|
| Init module              | `go mod init github.com/user/app`                               |
| Add dependency           | `go get github.com/pkg@v1.2.3`                                  |
| Clean deps               | `go mod tidy`                                                   |
| Vendor deps              | `go mod vendor`                                                 |
| Build                    | `go build -o app ./cmd/app/`                                    |
| Build (production)       | `go build -ldflags "-s -w" -o app ./cmd/app/`                   |
| Cross-compile            | `GOOS=linux GOARCH=amd64 go build -o app .`                     |
| Run                      | `go run ./cmd/app/`                                             |
| Test (all)               | `go test -race -cover ./...`                                    |
| Escape analysis          | `go build -gcflags="-m" ./...`                                  |
| Docker build             | `docker build -t app:latest .`                                  |
| Docker run               | `docker run -p 8080:8080 -e APP_ENV=prod app:latest`            |
| Compose up               | `docker compose up --build -d`                                  |
| Compose down             | `docker compose down -v`                                        |
| Push to registry         | `docker tag app:latest user/app:latest && docker push user/app:latest` |
| Lint                     | `golangci-lint run ./...`                                       |
| Install CLI tool         | `go install github.com/tool@latest`                             |

---

## Related Chapters

- [Chapter 13 — Memory, GC & Escape Analysis](./13_memory_gc_escape_sorting.md) — build flags for escape analysis (`-gcflags="-m"`)
- [Chapter 15 — Debugging & Profiling](./15_debugging_profiling.md) — pprof, dlv, GODEBUG
- [Chapter 18 — Production Patterns](./18_production_patterns_enterprise.md) — graceful shutdown, middleware, observability
