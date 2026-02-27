// Package build_deploy covers building, running, and deploying Go applications.
//
// ============================================================
// BUILD & RUN — COMPLETE GUIDE
// ============================================================
//
// ─────────────────────────────────────────────────────────────
// RUN WITHOUT BUILDING (dev mode)
// ─────────────────────────────────────────────────────────────
//   go run main.go                   ← run a single file
//   go run .                         ← run all .go files in current dir
//   go run ./cmd/myapp/              ← run a specific cmd package
//
// ─────────────────────────────────────────────────────────────
// BUILD — compile to a binary
// ─────────────────────────────────────────────────────────────
//   go build .                       ← builds binary in current dir
//   go build -o myapp .              ← name the output binary
//   go build ./...                   ← build all packages (no output file)
//   go build ./cmd/myapp/            ← build a specific package
//
// ─────────────────────────────────────────────────────────────
// BUILD FLAGS
// ─────────────────────────────────────────────────────────────
//   -o <name>         output binary name
//   -v                verbose (show packages being compiled)
//   -race             enable race detector (use in testing!)
//   -gcflags="-N -l"  disable optimisations (for debugger)
//   -ldflags          linker flags (see below for version injection)
//   -tags             build tags (e.g. -tags integration)
//
// ─────────────────────────────────────────────────────────────
// INJECT VERSION AT BUILD TIME (common production pattern)
// ─────────────────────────────────────────────────────────────
//   In your code (main.go):
//     var (
//         version   = "dev"
//         buildTime = "unknown"
//         gitCommit = "unknown"
//     )
//
//   Build command:
//     go build \
//       -ldflags "-X main.version=1.0.0 \
//                 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
//                 -X main.gitCommit=$(git rev-parse --short HEAD)" \
//       -o myapp .
//
// ─────────────────────────────────────────────────────────────
// CROSS-COMPILATION — build for a different OS/ARCH
// ─────────────────────────────────────────────────────────────
//   GOOS=linux   GOARCH=amd64 go build -o myapp-linux .
//   GOOS=windows GOARCH=amd64 go build -o myapp.exe .
//   GOOS=darwin  GOARCH=arm64 go build -o myapp-mac .
//
//   List all supported targets:
//     go tool dist list
//
//   Common GOOS:   linux, windows, darwin, freebsd
//   Common GOARCH: amd64, arm64, 386, arm
//
// ─────────────────────────────────────────────────────────────
// STRIP DEBUG INFO (smaller binary)
// ─────────────────────────────────────────────────────────────
//   go build -ldflags "-s -w" -o myapp .
//   -s: omit symbol table
//   -w: omit DWARF debug info
//   Result: ~30% smaller binary
//
// ─────────────────────────────────────────────────────────────
// INSTALL — build + place binary in $GOPATH/bin
// ─────────────────────────────────────────────────────────────
//   go install .
//   go install github.com/some/tool@latest   ← install a CLI tool
//
// ─────────────────────────────────────────────────────────────
// TYPICAL PROJECT STRUCTURE FOR A DEPLOYABLE APP
// ─────────────────────────────────────────────────────────────
//   myapp/
//   ├── cmd/
//   │   └── myapp/
//   │       └── main.go       ← entry point (package main)
//   ├── internal/             ← private packages (not importable externally)
//   │   ├── config/
//   │   ├── handler/
//   │   └── service/
//   ├── pkg/                  ← public reusable packages
//   ├── configs/
//   │   ├── config.yaml
//   │   └── config.local.yaml
//   ├── scripts/
//   │   └── build.sh
//   ├── Dockerfile
//   ├── docker-compose.yml
//   ├── Makefile
//   ├── go.mod
//   └── go.sum
//
// ─────────────────────────────────────────────────────────────
// MAKEFILE (common in Go projects)
// ─────────────────────────────────────────────────────────────
//   .PHONY: build run test clean
//
//   build:
//       go build -ldflags "-s -w" -o bin/myapp ./cmd/myapp/
//
//   run:
//       go run ./cmd/myapp/
//
//   test:
//       go test -race -cover ./...
//
//   clean:
//       rm -rf bin/
//
//   lint:
//       golangci-lint run ./...
//
// ─────────────────────────────────────────────────────────────
// RUNNING A GO APPLICATION
// ─────────────────────────────────────────────────────────────
//   Development:   go run .
//   After build:   ./myapp
//   With args:     ./myapp --port=8080 --env=production
//   With env vars: PORT=8080 ENV=production ./myapp
//
// ─────────────────────────────────────────────────────────────
// ENVIRONMENT VARIABLES IN THE SHELL
// ─────────────────────────────────────────────────────────────
//   # Linux / macOS
//   export APP_PORT=8080
//   export DB_URL=postgres://localhost/mydb
//   ./myapp
//
//   # Windows PowerShell
//   $env:APP_PORT="8080"
//   $env:DB_URL="postgres://localhost/mydb"
//   .\myapp.exe
//
//   # Inline (Linux / macOS)
//   APP_PORT=8080 DB_URL=postgres://localhost/mydb ./myapp
//
// ─────────────────────────────────────────────────────────────
// SYSTEMD SERVICE (Linux deployment)
// ─────────────────────────────────────────────────────────────
//   # /etc/systemd/system/myapp.service
//   [Unit]
//   Description=My Go App
//   After=network.target
//
//   [Service]
//   Type=simple
//   User=appuser
//   WorkingDirectory=/opt/myapp
//   EnvironmentFile=/opt/myapp/.env
//   ExecStart=/opt/myapp/myapp
//   Restart=on-failure
//   RestartSec=5s
//
//   [Install]
//   WantedBy=multi-user.target
//
//   # Commands:
//   systemctl enable myapp
//   systemctl start  myapp
//   systemctl status myapp
//   journalctl -u myapp -f   ← tail logs

package build_deploy
