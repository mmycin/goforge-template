# GoForge Template

This is the default project template for [GoForge](https://github.com/mmycin/GoForge), a production-ready Go framework and CLI.

## Prerequisites
- Go 1.21 or later
- GoForge CLI

## Getting Started
This template is automatically used when you run `goforge new <project-name>`. For full documentation, please visit the [GoForge Documentation](https://github.com/mmycin/GoForge/tree/main/goforge-cli/docs).

## Project Structure
```text
your-project/
├── boot/               # Core framework code (not modified often)
│   ├── kernel.go       # Application bootstrap
│   ├── client/         # gRPC and HTTP client wrappers
│   └── server/         # HTTP/gRPC server implementations
├── cmd/                # Entry point
├── internal/           # Your application code
│   ├── config/         # Configuration
│   ├── console/        # Custom commands
│   ├── database/       # Database migrations/queries
│   └── services/       # Domain-specific services
├── proto/              # gRPC definitions
├── tests/              # Test files
├── .env.example        # Env file template
├── Dockerfile          # Docker config
├── air.toml            # Live-reload config
└── sqlc.yaml           # SQLC config
```

## License
GoForge is licensed under the Apache License 2.0. See [LICENSE](LICENSE).
