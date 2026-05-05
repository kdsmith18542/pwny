# Pwny Architecture

## Design Principles

1. **API-First** — Everything the GUI does goes through REST + WebSocket
2. **Headless by Default** — The engine runs as a server; GUI and CLI are equal consumers
3. **Static Module Registration** — No Go plugins; modules register via `init()` at compile time
4. **Event-Driven** — Internal event bus pushes state changes to WebSocket clients
5. **Graceful Degradation** — Operates in memory-only mode if SQLite is unavailable

## Layers

```
┌──────────────────────────────────────┐
│           Tauri Desktop App           │
│     (React + xterm.js + Zustand)      │
└──────────────┬───────────────────────┘
               │ HTTP/WS (localhost:31337)
┌──────────────┴───────────────────────┐
│         HTTP API Server (Go)          │
│  chi router + gorilla/websocket       │
│  ┌─────────┐ ┌─────────┐ ┌────────┐  │
│  │ Modules │ │Session  │ │Payload │  │
│  └────┬────┘ └────┬────┘ └───┬────┘  │
│       │           │          │        │
│  ┌────┴────┐ ┌────┴────┐ ┌──┴─────┐  │
│  │ Registry│ │Session  │ │DB Layer│  │
│  │         │ │Manager  │ │(SQLite)│  │
│  └─────────┘ └─────────┘ └────────┘  │
└──────────────────────────────────────┘
```

## Core Components

### Module System (internal/core/)

- **Module interface**: `Info()`, `Options()`, `SetOption()`, `Validate()`, `Run()`
- **BaseModule**: Default implementation with option tracking
- **Registry**: Static `map[string]ModuleFactory` with init()-time registration
- **Config**: YAML via viper, environment variable overrides (`PWNY_*`)

### Session Management (internal/api/ + internal/core/)

- **Session interface**: `Read`, `Write`, `Execute`, `Close`, `Upload`, `Download`
- **SessionManager**: Thread-safe registry backed by map, create/list/get/close
- **WebSocket relay**: Real-time session I/O via gorilla/websocket

### Persistence (internal/db/)

- **SQLite**: Workspaces, hosts, services, credentials, notes, event_log
- **BoltDB**: Reserved for ephemeral session state (future)

### API Server (internal/api/)

- **HTTP**: chi router with CORS, logging, recovery middleware
- **WebSocket**: Session terminal I/O, event stream for push updates
- **Endpoints**: Module CRUD, Session CRUD, status, events

## Database Schema

Tables: `workspaces`, `hosts`, `services`, `credentials`, `notes`, `loots`, `event_log`

See [internal/db/database.go](internal/db/database.go) for current DDL.

## Module System

Modules are Go source files in `modules/` that self-register via `init()`:

```go
func init() {
    core.Register("exploit/windows/smb/ms17_010", func() core.Module {
        return &MS17_010{...}
    })
}
```

This replaces the prior Go-plugin approach (Linux-only, version-sensitive, removed).

## Future

- **GUI**: Tauri + React desktop app
- **Payload engine**: Stager/stage generation with encoders
- **Module library**: Curated top-50 exploits and auxiliary scanners
- **Event bus**: Pub/sub for cross-component notifications
