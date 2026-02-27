// Package docker demonstrates how to Dockerize a Go application.
//
// ============================================================
// DOCKERIZING A GO APPLICATION — COMPLETE GUIDE
// ============================================================
//
// ─────────────────────────────────────────────────────────────
// SIMPLE DOCKERFILE (single-stage, good for learning)
// ─────────────────────────────────────────────────────────────
//
//   FROM golang:1.25-alpine
//
//   WORKDIR /app
//
//   # Copy go.mod and go.sum first (layer caching — only re-download if these change)
//   COPY go.mod go.sum ./
//   RUN go mod download
//
//   # Copy the rest of the source code
//   COPY . .
//
//   # Build the binary
//   RUN go build -o myapp ./cmd/myapp/
//
//   EXPOSE 8080
//
//   CMD ["./myapp"]
//
// ─────────────────────────────────────────────────────────────
// PRODUCTION DOCKERFILE — multi-stage build (RECOMMENDED)
// ─────────────────────────────────────────────────────────────
//   Multi-stage builds give you:
//   ✅ Final image has NO Go toolchain (much smaller)
//   ✅ Source code NOT included in final image (security)
//   ✅ Typical result: 10-15 MB instead of 300+ MB
//
//   ─── Dockerfile ───────────────────────────────────────────
//
//   # ── Stage 1: Build ───────────────────────────────────────
//   FROM golang:1.25-alpine AS builder
//
//   WORKDIR /app
//
//   # Install git if private repos are needed
//   RUN apk add --no-cache git
//
//   COPY go.mod go.sum ./
//   RUN go mod download
//
//   COPY . .
//
//   # CGO_ENABLED=0: static binary (no libc dependency)
//   # -ldflags "-s -w": strip debug info → smaller binary
//   RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
//       go build -ldflags="-s -w" -o myapp ./cmd/myapp/
//
//   # ── Stage 2: Final image ─────────────────────────────────
//   FROM scratch
//   # OR use: FROM gcr.io/distroless/static:nonroot
//   # OR use: FROM alpine:3.19  (if you need shell/certs/etc.)
//
//   # Copy CA certificates (needed for HTTPS calls)
//   COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
//
//   # Copy the binary
//   COPY --from=builder /app/myapp /myapp
//
//   # Copy config files if needed
//   COPY --from=builder /app/configs/ /configs/
//
//   # Run as non-root (security best practice)
//   USER 65534:65534
//
//   EXPOSE 8080
//
//   ENTRYPOINT ["/myapp"]
//
// ─────────────────────────────────────────────────────────────
// .dockerignore  (put this next to Dockerfile)
// ─────────────────────────────────────────────────────────────
//   .git
//   .gitignore
//   *.md
//   bin/
//   vendor/
//   .env
//   .env.local
//   *_test.go
//
// ─────────────────────────────────────────────────────────────
// BUILD & RUN DOCKER COMMANDS
// ─────────────────────────────────────────────────────────────
//   # Build image
//   docker build -t myapp:latest .
//   docker build -t myapp:1.0.0 .
//   docker build --no-cache -t myapp:latest .   ← force fresh build
//
//   # Run container
//   docker run myapp:latest
//   docker run -p 8080:8080 myapp:latest          ← map host:container ports
//   docker run -d -p 8080:8080 myapp:latest       ← detached (background)
//   docker run --name my-running-app -p 8080:8080 myapp:latest
//
//   # Pass environment variables
//   docker run -e APP_ENV=production -e DB_URL=postgres://... myapp:latest
//   docker run --env-file .env myapp:latest       ← load from .env file
//
//   # Mount a volume (for config files)
//   docker run -v $(pwd)/configs:/configs myapp:latest
//
//   # View logs
//   docker logs <container-id>
//   docker logs -f my-running-app    ← follow/tail logs
//
//   # Stop / remove
//   docker stop my-running-app
//   docker rm   my-running-app
//
//   # List running containers
//   docker ps
//   docker ps -a   ← all including stopped
//
//   # Execute command inside running container
//   docker exec -it my-running-app sh
//
// ─────────────────────────────────────────────────────────────
// DOCKER COMPOSE — for multi-service local dev
// ─────────────────────────────────────────────────────────────
//   ─── docker-compose.yml ───────────────────────────────────
//
//   version: "3.9"
//
//   services:
//     app:
//       build: .
//       ports:
//         - "8080:8080"
//       environment:
//         - APP_ENV=development
//         - DB_URL=postgres://postgres:password@db:5432/mydb
//       depends_on:
//         db:
//           condition: service_healthy
//       volumes:
//         - ./configs:/configs
//
//     db:
//       image: postgres:16-alpine
//       environment:
//         POSTGRES_USER: postgres
//         POSTGRES_PASSWORD: password
//         POSTGRES_DB: mydb
//       ports:
//         - "5432:5432"
//       healthcheck:
//         test: ["CMD-SHELL", "pg_isready -U postgres"]
//         interval: 5s
//         timeout: 5s
//         retries: 5
//       volumes:
//         - pgdata:/var/lib/postgresql/data
//
//   volumes:
//     pgdata:
//
//   ──────────────────────────────────────────────────────────
//
//   # Docker Compose commands:
//   docker compose up              ← start all services
//   docker compose up -d           ← start in background
//   docker compose up --build      ← rebuild images first
//   docker compose down            ← stop and remove containers
//   docker compose down -v         ← also remove volumes
//   docker compose logs -f app     ← follow logs for 'app' service
//   docker compose ps              ← list service status
//   docker compose exec app sh     ← shell into running app container
//
// ─────────────────────────────────────────────────────────────
// PUSH TO REGISTRY (Docker Hub / AWS ECR / GCP Artifact Registry)
// ─────────────────────────────────────────────────────────────
//   # Docker Hub
//   docker login
//   docker tag myapp:latest yourusername/myapp:latest
//   docker push yourusername/myapp:latest
//
//   # AWS ECR
//   aws ecr get-login-password --region us-east-1 | \
//     docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
//   docker tag myapp:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/myapp:latest
//   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/myapp:latest
//
// ─────────────────────────────────────────────────────────────
// HOT RELOAD IN DOCKER (development only)
// ─────────────────────────────────────────────────────────────
//   Use 'air' for live reload:
//   go install github.com/air-verse/air@latest
//
//   ─── Dockerfile.dev ───────────────────────────────────────
//   FROM golang:1.25-alpine
//   RUN go install github.com/air-verse/air@latest
//   WORKDIR /app
//   COPY go.mod go.sum ./
//   RUN go mod download
//   CMD ["air"]
//   ──────────────────────────────────────────────────────────
//
//   Then mount your source:
//   docker run -v $(pwd):/app -p 8080:8080 myapp-dev

package docker
