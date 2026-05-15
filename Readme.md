# GoForge Framework

![GoForge Banner](https://raw.githubusercontent.com/mmycin/GoForge/main/assets/logo_without_bg.png)

> [!IMPORTANT]
> A comprehensive, production-ready Go application framework designed for high-performance gRPC/HTTP services, robust database management, and rapid scaffolding.

GoForge provides a powerful foundation for building scalable Go applications. It eliminates boilerplate by providing integrated tooling for database migrations, type-safe SQL, and gRPC service generation.

---

## ✨ Features

- **🚀 Unified CLI (`goforge`)**: A single tool to manage your entire development lifecycle—from scaffolding to migrations.
- **🏗️ Service Scaffolding**: Generate complete domain layers (services, repositories, gRPC stubs, and HTTP handlers) in seconds.
- **🗄️ Database Excellence**: Native support for **Atlas** (migrations) and **SQLC** (type-safe queries) across SQLite, MySQL, and PostgreSQL.
- **📡 Dual-Protocol Servers**: Run high-performance HTTP (Fiber) and gRPC servers concurrently with shared middleware and context.
- **🛠️ Production Ready**: Structured JSON logging, environment configuration, rate limiting, and secure application key encryption.

---

## 🏗️ Project Structure

GoForge follows a clean, modular architecture optimized for separation of concerns:

```text
.
├── cmd/                # Entry points (main.go)
├── internal/           # Private application code
│   ├── client/         # Internal gRPC/HTTP client wrappers
│   ├── config/         # Environment & config loaders
│   ├── console/        # Custom application CLI commands
│   ├── database/       # Migrations, SQLC gen, and DB core
│   ├── server/         # HTTP/gRPC server implementations
│   └── services/       # Domain business logic & models
├── proto/              # Protobuf definitions
├── tests/              # Comprehensive test suites
├── air.toml            # Live-reloading config
├── atlas.hcl           # Atlas migration config
└── sqlc.yaml           # SQLC query generation config
```

---

## 🏁 Getting Started

### 1. Install the GoForge CLI

Install the universal CLI tool to your `$GOPATH/bin`:

```bash
go install github.com/mmycin/GoForge@latest
```

### 2. Environment Setup

Initialize your environment by copying the example and generating a secure application key:

```bash
cp .env.example .env
goforge gen:key
```

### 3. Initialize & Run

Tidy your dependencies and start the development servers:

```bash
go mod tidy
goforge app serve
```

---

## 🛠️ Command Line Interface (CLI)

The `goforge` CLI is your primary interface for development.

### 🔑 Core Commands

| Command | Description |
| :--- | :--- |
| `goforge gen:key` | Generates a 32-character `APP_KEY` in your `.env`. |
| `goforge app serve` | Starts the application (HTTP & gRPC servers). |
| `goforge version` | Checks your current GoForge CLI version. |

### 🏗️ Scaffolding & Services

Accelerate your development by generating boilerplate-free services:

- **`goforge gen:service [name]`**: Creates a full service stack in `internal/services/`.
- **`goforge rem:service [name]`**: Safely removes a service and its associated registrations.
- **`goforge gen:command [name]`**: Bootstraps a new console command in `internal/console/`.

### 🗄️ Database & Migrations

GoForge uses **Atlas** for declarative migrations and **SQLC** for type-safe code generation.

- **`goforge gen:migration [name]`**: Generates a new migration by diffing GORM models against the DB.
- **`goforge migrate`**: Applies all pending migrations to your database.
- **`goforge rem:migration`**: Reverts the most recent migration file.
- **`goforge gen:sqlc`**: Compiles raw SQL in `internal/database/queries` into Go code.
- **`goforge loader`**: Displays the current GORM schema as interpreted by Atlas.

### 📡 Protocol Buffers (gRPC)

- **`goforge gen:proto [name]`**: Compiles `.proto` files into Go gRPC stubs.
- **`goforge rem:proto`**: Cleans up all generated `.pb.go` files.

---

## 🚀 Deployment

1. **Statically Linked Binary**: Build for Linux/Docker:
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
   ```

### 🐳 Docker Compose

You can spin up the entire stack (App + Postgres) using Docker Compose:

```bash
docker-compose up --build
```

This will:
1. Build the Go application using the multi-stage `Dockerfile`.
2. Start a PostgreSQL 16 container.
3. Automatically link them via a private network.
2. **Production Settings**:
   - Set `APP_DEBUG=false`
   - Set `LOG_FORMAT=json`
   - Ensure `APP_KEY` is set via environment secrets, not `.env` files.

---

## 📄 License

GoForge is open-sourced software licensed under the [Apache License 2.0](LICENSE).
