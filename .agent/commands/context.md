# Context

Use this command when starting work in `goforge-template`.

Read in this order:
1. `README.md`
2. `.agent/rules/architecture.md`
3. `.agent/rules/testing.md`
4. The service or package you will touch

Project facts:
- This repository is the generated application template, not the CLI product.
- Core framework code lives under `boot/`.
- Application code lives under `internal/`.
- Services are wired in `internal/services/kernel.go`.
- Avoid changing `goforge-cli` from this repo unless explicitly requested.

Non-negotiables:
- No package-level state in `internal/`.
- Use constructor-based dependency injection.
- Imports must use `github.com/mmycin/goforge/boot/...`.
- Keep HTTP, gRPC, config, and database changes consistent with existing patterns.
