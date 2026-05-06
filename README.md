# Pwny

A modern, modular, cross-platform penetration testing framework with an HTTP API and Web UI (Tauri).

## Features

- **Modular Architecture**: Static registry pattern for compile-time-safe module loading
- **REST + WebSocket API**: Headless server mode, build any client on top
- **Session Management**: Shell and Meterpreter session support with real-time I/O via WebSocket
- **SQLite Persistence**: Workspace management, host/service tracking, credential store
- **Cross-Platform**: Single static binary for Windows, Linux, macOS
- **Configurable**: YAML config file, environment variable overrides

## Quick Start

```bash
# Prerequisites: Go 1.23+
git clone https://github.com/kdsmith18542/pwny.git
cd pwny

go mod tidy
go build -o bin/pwny-server ./cmd/pwny-server/

# Generate default config, then edit pwny.yaml
./bin/pwny-server init-config

# Start the API server
./bin/pwny-server -c pwny.yaml
```

## Project Structure

```
pwny/
├── cmd/pwny-server/     # API server entry point
├── internal/
│   ├── core/            # Framework kernel (module interface, registry, config, errors)
│   ├── session/         # Session management (reserved)
│   ├── payload/         # Payload generation (reserved)
│   ├── api/             # REST + WebSocket handlers
│   ├── db/              # SQLite persistence + migrations
│   └── network/         # Network primitives (reserved)
├── modules/             # Module library (exploit, auxiliary, post)
├── examples/            # Example modules
├── gui/                 # Tauri desktop app (planned)
├── docs/                # Documentation
├── Taskfile.yml         # Build automation
└── .github/workflows/   # CI pipeline
```

## API

Server runs on `http://127.0.0.1:31337` by default.

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/status` | Server health |
| `GET` | `/api/v1/modules` | List modules (`?type=exploit` filter) |
| `GET` | `/api/v1/modules/{name}` | Module details + options |
| `POST` | `/api/v1/modules/{name}/validate` | Validate module options |
| `POST` | `/api/v1/modules/{name}/run` | Execute a module |
| `GET` | `/api/v1/sessions` | List sessions |
| `GET` | `/api/v1/sessions/{id}` | Session details |
| `DELETE` | `/api/v1/sessions/{id}` | Close a session |
| `GET` | `/api/v1/sessions/{id}/ws` | WebSocket session I/O |
| `GET` | `/api/v1/events` | WebSocket event stream |

## Developing

```bash
# Run all tests
go test -count=1 ./...

# Run with verbose output
go test -v -count=1 ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code
gofmt -l -w .

# Build everything
go build ./...
```

See [Taskfile.yml](Taskfile.yml) for more targets (`task build`, `task test`, `task lint`, etc.).

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md) — to be written
- [Module Development](docs/MODULE_DEVELOPMENT.md) — to be written

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
