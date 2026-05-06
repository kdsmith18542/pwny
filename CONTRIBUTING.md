# Contributing to Pwny

## Getting Started

1. Fork the repository
2. Clone your fork
3. Run `go mod tidy`

## Development Workflow

1. Create a branch: `git checkout -b feature/your-feature`
2. Make changes, following existing code style
3. Run tests: `go test -count=1 ./...`
4. Format: `gofmt -l -w .`
5. Lint: `go vet ./...`
6. Commit with a descriptive message
7. Push and open a PR against `main`

## Testing

- Write tests for new features and bug fixes
- Run all tests: `go test -count=1 ./...`
- Coverage: `go test -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out`

## Documentation

When adding new features:
1. Update relevant `.md` files
2. Add GoDoc comments for exported types and functions
3. Update progress in `docs/PROGRESS.md`

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Meaningful variable and function names
- Write unit tests for new functionality

## Module Development

Modules register themselves via `init()`:

```go
package mymodule

import "github.com/kdsmith18542/pwny/internal/core"

type MyModule struct {
    *core.BaseModule
}

func New() *MyModule {
    m := &MyModule{BaseModule: core.NewBaseModule(core.TypeAuxiliary, "my_module")}
    m.RegisterOption("RHOST", "Target host", true, nil)
    return m
}

func (m *MyModule) Run() (interface{}, error) {
    return "module output", nil
}

func init() {
    core.Register("auxiliary/my_module", func() core.Module {
        return New()
    })
}
```
