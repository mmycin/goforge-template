# TLS / HTTPS

Use this skill when a task touches TLS configuration or certificate management.

Read first:
- `internal/config/http.go` — `HTTPConfig.TLSEnabled`, `TLSCertFile`, `TLSKeyFile`
- `boot/server/http.go` — `Start()` branches on `TLSEnabled`

Configuration (`.env`):
```
TLS_ENABLED=true
TLS_CERT_FILE=cert.pem
TLS_KEY_FILE=key.pem
```

Local development (self-signed cert):
```bash
openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem -out cert.pem \
  -days 365 -nodes \
  -subj "/CN=localhost"
```

Production:
- Use certificates from your CA or a managed service (e.g., Let's Encrypt via `golang.org/x/crypto/acme/autocert`).
- Point `TLS_CERT_FILE` / `TLS_KEY_FILE` to the provisioned paths.

Guardrails:
- Never commit `.pem` or `.key` files to version control (they are already in `.gitignore`).
- `TLS_ENABLED=false` is the safe default for development; enable it explicitly in production.
- If TLS is enabled but cert/key paths are missing, the server refuses to start with a clear error.
