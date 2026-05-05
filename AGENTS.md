# Pwny — Agent Context

## Project Overview

Modern, modular, cross-platform penetration testing framework.
Go backend with REST + WebSocket API, Tauri desktop GUI (planned).

## Architecture

```
cmd/pwny-server/          # CLI entry point (Cobra)
internal/core/             # Framework kernel
  module.go                # Module interface, BaseModule, option system
  registry.go              # Static module registry (init()-time registration)
  config.go                # YAML config via viper, PWNY_* env vars
  errors.go                # Sentinel errors
  session.go               # Session interface, SessionManager, BaseSession
  shell_session.go         # Shell session implementation
  meterpreter_session.go   # Meterpreter session (TLV protocol, partial)
internal/api/              # HTTP API layer
  server.go                # Server struct, shared sessionManager
  router.go                # chi routes, APIResponse/APIError types, writeJSON/writeError
  middleware.go             # Logging, panic recovery
  handler_module.go        # Module list/get/validate/run
  handler_session.go       # Session list/get/close
  websocket.go             # Session I/O relay via gorilla/websocket
internal/db/               # Persistence
  database.go              # SQLite open + schema migration
  workspace.go             # Workspace CRUD
modules/                   # Module library (self-registering via init())
```

## Key Design Decisions

- **No Go plugins** — modules use static registry pattern (`core.Register()` in `init()`)
- **API-first** — server exposes REST + WebSocket; GUI/CLI are consumers
- **Headless** — server runs standalone; GUI is a separate process
- **SQLite** — WAL mode, no CGO (modernc.org/sqlite)

## Build & Test

```bash
go build ./...                         # build all
go build -o bin/pwny-server ./cmd/pwny-server/  # server binary
go test -count=1 ./...                 # all tests
go test -v -count=1 ./internal/core/   # core tests verbose
go vet ./...                           # lint
gofmt -l -w .                          # format
task build                             # via Taskfile.yml
task test
```

Module registration pattern:

```go
func init() {
    core.Register("exploit/windows/smb/ms17_010", func() core.Module {
        return &MS17_010{BaseModule: core.NewBaseModule(core.TypeExploit, "windows/smb/ms17_010")}
    })
}
```

## Dependencies

- github.com/go-chi/chi/v5 — HTTP router
- github.com/gorilla/websocket — WebSocket
- github.com/spf13/cobra — CLI framework
- github.com/spf13/viper — Config
- modernc.org/sqlite — SQLite (no CGO)
- go.etcd.io/bbolt — KV store (planned)
- github.com/stretchr/testify — Test assertions
- github.com/google/uuid — UUID generation

## API

Server defaults to `127.0.0.1:31337`.

| Method | Path | Description |
|---|---|---|
| GET | /api/v1/status | Health check |
| GET | /api/v1/modules | List modules (?type=exploit) |
| GET | /api/v1/modules/{name} | Module details + options |
| POST | /api/v1/modules/{name}/validate | Validate options |
| POST | /api/v1/modules/{name}/run | Execute module |
| GET | /api/v1/sessions | List sessions |
| GET | /api/v1/sessions/{id} | Session details |
| DELETE | /api/v1/sessions/{id} | Close session |
| GET | /api/v1/sessions/{id}/ws | WebSocket session I/O |
| GET | /api/v1/events | WebSocket event stream |

## Testing

- 40 unit tests in internal/core (100% module interface, 92% session manager, 100% registry)
- 44% statement coverage overall
- Co-located test files: `foo.go` → `foo_test.go`
- Use testify for assertions
- Mock network with `net.Pipe()`
- Table-driven tests for validation logic

## Commit Style

Imperative mood, short summary line, optional body with bullet points.

```
Implement Phase N: short description

- Specific change one
- Specific change two
- Specific change three
```
