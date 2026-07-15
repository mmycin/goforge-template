# API Design Rule

When changing HTTP or gRPC behavior:
- Keep transport thin; business logic belongs in services.
- Validate request input at the edge.
- Return consistent status codes and error shapes.
- Keep route registration local to the owning service package.
- For gRPC, maintain graceful shutdown behavior and avoid long-running uncontrolled goroutines.
- Prefer additive API changes unless the user explicitly requests a breaking change.
