# Architecture Rule

Treat this repository as a generated application built on GoForge.

Structure:
- `boot/` contains framework bootstrap and server infrastructure.
- `internal/` contains application-specific code.
- `internal/services/kernel.go` is the composition root for service wiring.

Rules:
- Do not introduce package-level mutable state in `internal/`.
- Prefer constructor injection with `New*` functions.
- Keep framework bootstrap concerns in `boot/`, not in service packages.
- Keep service-specific behavior inside its own `internal/services/<name>/` package.
- If a change belongs in the generator rather than the generated app, stop and confirm whether `goforge-cli` should be updated too.
