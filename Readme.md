# GoForge

GoForge is a minimal, lightweight, fast, and DX-friendly microservice starter template for Golang. Built with a modular architecture and powered by top-tier Go tools.

## 🚀 Key Features

- **CLI Interface**: Powered by Cobra
- **Configuration**: Structured and type-safe using Viper (`config.DB.Name`)
- **HTTP Server**: Fast and robust REST APIs with Gin
- **Database**:
    - **ORM**: GORM for easy relational mapping
    - **Query**: SQLC for type-safe, performant raw SQL
    - **Migrations**: Versioned migrations using Atlas
- **Developer Experience**:
    - Modular service structure
    - Automatic route registration
    - Graceful shutdown support
    - Service scaffolding via CLI

## 🛠️ Getting Started

### Prerequisites

- Go 1.24+
- SQLC (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- Atlas (`go install ariga.io/atlas/cmd/atlas@latest`)

### Initial Setup

1. Clone the repository
2. Install dependencies:
    ```bash
    go mod tidy
    ```
3. Configure your environment:
    ```bash
    cp .env.example .env # Ensure .env exists
    ```

## 💻 CLI Usage

GoForge provides a powerful CLI to speed up development.

### Start the API Server

```bash
go run . serve
```

### Authentication & Security

GoForge includes an `App-Key` middleware by default to protect your APIs.

1. **Generate a Key**:
    ```bash
    go run . gen:key
    ```
2. **Usage**:
   Include the key in your request headers:
    ```http
    X-App-Key: your_generated_key
    ```
    _Note: `/health` remains publicly accessible._

### Database Operations

```bash
# Generate a new migration
go run . gen:migration create_users_table

# Run pending migrations
go run . migrate

# Display database schema
go run . loader
```

### Code Generation

```bash
# Generate a new modular service
go run . make:service user

# Generate SQLC query code
go run . gen:sqlc
```

## 🏗️ Architecture

GoForge follows a modular architecture where each feature is encapsulated within its own service directory.

```
internal/services/todo/
  ├── handler.go  # HTTP request handlers (Gin)
  ├── routes.go   # Route definitions & registration
  ├── model.go    # GORM models
  ├── repo.go     # Repository layer (SQLC/GORM)
  ├── service.go  # Business logic layer
  └── proto.go    # gRPC definitions
```

### Route Registration

When you create a new service, it's automatically registered in the application router. GoForge provides a clean, callback-based grouping DX:

```go
// routes.go
func (r *TodoRouter) Register(engine gin.IRouter) {
    h := &TodoHandler{}

    // DX-friendly grouping helper
    RegisterGroup(engine, "/todos", func(group *gin.RouterGroup) {
        group.GET("/", h.GetAllTodos)
        group.GET("/:id", h.GetTodoByID)
    })
}
```

---

## ⚙️ Configuration

GoForge uses **Viper** for type-safe configuration. Settings are grouped by component and accessible via the `config` package.

### Usage

```go
import "github.com/mmycin/goforge/internal/config"

dbName := config.DB.Name
port := config.App.Port
```

### Adding New Config

1. Add your variable to `.env`.
2. Define the field in the corresponding struct in `internal/config/`.
3. Viper will automatically map the environment variable to your struct during `config.Load()`.

---

## 🛑 Graceful Shutdown

GoForge listens for termination signals (`SIGINT`, `SIGTERM`) and allows the server up to 30 seconds to finish processing active requests before shutting down. This ensures zero-downtime deployments and data integrity.

---

## 🧪 Testing

GoForge encourages Test-Driven Development (TDD).

```bash
# Run all tests
go test ./...

# Run service-specific tests
go test ./internal/services/todo/...
```

---

## 🔮 gRPC Support

GoForge supports running a gRPC server concurrently with the HTTP server.

1.  **Enable gRPC**: Set `GRPC_ENABLE=true` in `.env`.
2.  **Define Service**: Create your `.proto` file in `proto/<service>/<service>.proto`.
3.  **Generate Code**:
    ```bash
    go run . gen:proto <service_name>
    ```
4.  **Implement Server**: Edit `internal/services/<service>/grpc.go` to implement your logic.

---

## 📜 License

MIT
