# Testing

Use this skill when a task changes runtime behavior or shared wiring.

Test map:
- `tests/http/` for HTTP-facing behavior
- `tests/grpc/` for gRPC-facing behavior
- `tests/integration/` for cross-layer behavior
- `tests/unit/` for smaller package-level checks

Workflow:
1. Run the narrowest useful tests first.
2. Expand to broader coverage if shared code changed.
3. Record what was verified and what remains unverified.

Guardrails:
- Do not claim verification you did not run.
- Missing tests are a real finding, not a footnote.
