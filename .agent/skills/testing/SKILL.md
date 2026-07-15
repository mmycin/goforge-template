# Testing

Use this skill when a task changes runtime behavior or shared wiring.

Test map:
- `internal/tests/http/` for HTTP-facing behavior
- `internal/tests/grpc/` for gRPC-facing behavior
- `internal/tests/integration/` for cross-layer behavior
- `internal/tests/unit/` for smaller package-level checks

Scripts:
```bash
bash .agent/skills/testing/run.sh          # run all scopes
bash .agent/skills/testing/run.sh http     # HTTP only
bash .agent/skills/testing/run.sh grpc     # gRPC only
bash .agent/skills/testing/run.sh integration
bash .agent/skills/testing/run.sh unit
```

Workflow:
1. Run the narrowest useful tests first.
2. Expand to broader coverage if shared code changed.
3. Record what was verified and what remains unverified.

Guardrails:
- Do not claim verification you did not run.
- Missing tests are a real finding, not a footnote.
