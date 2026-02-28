# GoForge Framework

GoForge is a comprehensive, production-ready Go application framework designed to provide robust tooling for database migrations, code generation, caching, and service scaffolding. It comes with built-in support for multiple SQL databases, gRPC and HTTP servers, tiered caching, and a powerful CLI to speed up development.

---

## 🚀 Features

- **Built-in CLI (`goforge`)**: Scaffolding for services, database migrations, gRPC protobufs, and more.
- **Multiple Databases**: Out-of-the-box support for SQLite, MySQL, PostgreSQL, and SQL Server.
- **Advanced Caching**: Scalable caching layer supporting in-memory (Ristretto), distributed (Redis), or a tiered combination of both.
- **Code Generation**: Automated integration for **SQLC** (type-safe SQL) and **protoc** (gRPC).
- **Dual Server Support**: Run HTTP and gRPC servers concurrently with ease.
- **Production Ready**: Includes rate limiting, structured logging, encryption capabilities, and environment-driven configuration.

---

## 📦 Getting Started

### 1. Installation

Clone your repository and install dependencies:

```bash
go mod tidy
```

### 2. Environment Configuration

Copy the `.env.example` file to create your local `.env`:

```bash
cp .env.example .env
```

Generate a secure application key:

```bash
go run cmd/main.go gen:key
```

Configure your database and caching options within the newly created `.env` file (see details below).

### 3. Running the Server

Start both the HTTP and gRPC servers:

```bash
go run cmd/main.go serve
```

---

## 🛠 Command Line Interface (CLI)

GoForge provides an extensive CLI to automate repetitive developer tasks.

### Service Scaffolding

- `go run cmd/main.go gen:service <name>`: Generates a complete service footprint (HTTP routes, controllers, repository layers).
- `go run cmd/main.go rem:service <name>`: Removes a specified service.

### Database & Migrations

- `go run cmd/main.go gen:migration <name>`: Creates a new timestamped migration file.
- `go run cmd/main.go migrate`: Executes pending database migrations.
- `go run cmd/main.go rem:migration`: Reverts the latest database migration.
- `go run cmd/main.go loader`: Runs the GORM schema loader.

### Code Generation Integrations

- `go run cmd/main.go gen:sqlc`: Initializes and generates type-safe Go code from your raw SQL queries (using `sqlc`).
- `go run cmd/main.go rem:sqlc`: Removes SQLC integration, instantly wiping the generated models to keep the repo clean.
- `go run cmd/main.go gen:proto`: Compiles your `.proto` files into Go gRPC/protobuf code.
- `go run cmd/main.go rem:proto`: Cleans up generated `.proto` code.

---

## 🗄️ Database Setup

GoForge manages databases centrally. Supported database connection types (`DB_CONNECTION` in `.env`):

- `sqlite`
- `mysql`
- `postgres`
- `sqlserver`

You can use the built-in migration system (configurable to use `atlas` or `gorm` via `DB_MIGRATOR`) to safely manage schema changes over time.

---

## ⚡ Caching Layer

The flexible caching layer avoids deep refactoring as your scalable needs grow. The configured driver is universally accessible via `cache.Global`.

### Configuration (`.env`)

```ini
CACHE_ENABLED=true
# Options: memory | redis | both
CACHE_DRIVER=both
CACHE_TTL=5m
CACHE_MAX_ITEMS=10000
CACHE_MAX_COST=100MB

# Required if driver is redis or both
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Supported Drivers

1. **`memory`**: Uses **Ristretto** for blazing-fast, concurrent local caching. Best for single-instance, lightweight deployments.
2. **`redis`**: Connects to a standard Redis instance. Necessary for load-balanced environments with multiple instances requiring synchronized state.
3. **`both` (Tiered)**: A local L1 (Memory) and distributed L2 (Redis) approach. Read misses check Redis, then populate the local cache. Writes go to both. Optimizes extreme high-read scenarios while maintaining distributed consistency.

### Usage in Code

```go
import "github.com/mmycin/goforge/internal/cache"

// Set a cache value
cache.Set(ctx, "user:1", userObj, 15 * time.Minute)

// Retrieve a cache value (pass a pointer to the target struct)
var user User
err := cache.Get(ctx, "user:1", &user)

// Direct driver access (if absolutely necessary)
cache.Memory.Set(ctx, "local_flag", true, 0)
cache.Redis.Delete(ctx, "remote_key")
```

---

## 🚀 Deployment & Production Readiness

When taking your GoForge application to production, strictly adhere to these practices:

1. **Security**
    - **Environment Variables**: Never commit the `.env` file containing secrets `APP_KEY`, `DB_PASSWORD`, or `REDIS_PASSWORD`. Pass secrets through securely managed CI/CD pipelines or Secret Managers (AWS Secrets Manager, HashiCorp Vault).
    - **`APP_DEBUG=false`**: Always disable application debugging to prevent stack traces from leaking to the end user.

2. **Scaling & Caching**
    - **Multi-Instance Scaling**: If deploying multiple containers (e.g., Kubernetes, Docker Swarm), switch `CACHE_DRIVER` to `redis` or `both`. Using `memory` natively in a distributed setup will lead to cache inconsistency.
3. **Logging (`LOG_TYPE`, `LOG_FORMAT`)**
    - Output logs as JSON (`LOG_FORMAT=json`) rather than plain text for optimal parsing by observability tools (Datadog, ELK, Splunk).
    - Use `LOG_TYPE=both` or `LOG_TYPE=file` paired with a reliable log rotation strategy on your production servers.

4. **Resource Control**
    - **Rate Limiting**: Control traffic via `RATE_LIMIT_PER_MINUTE`. Adjust this upward for internal microservices, and downward for public facing endpoints to mitigate brute force/DDoS attacks.

5. **Compilation**
    - Ensure you are building statically linked Go binaries for production execution:
        ```bash
        CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
        ```
