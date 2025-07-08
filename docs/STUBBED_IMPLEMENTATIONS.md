# Stubbed Implementations Guide

## Purpose
This document defines the standards for creating and documenting stubbed implementations in the MSF-Go codebase.

## Stub Documentation Requirements

### 1. Stub Header
Every stubbed function/method must include a header comment with:

```go
// STUB: [Component] - [Function Name]
// Status: [Status]
// Priority: [P0-P4]
// Dependencies: [List of dependencies]
// Blockers: [List of blocking issues]
// Last Updated: YYYY-MM-DD
//
// Description: [Brief description of what this should do]
//
// TODO: [Specific tasks needed to complete]
// NOTE: [Any important notes]
// WARNING: [Any potential issues or limitations]
```

### 2. Status Levels
- `Not Started`: No work has begun
- `In Progress`: Actively being implemented
- `Partially Implemented`: Basic functionality exists but incomplete
- `Needs Testing`: Implementation complete but requires testing
- `Deprecated`: No longer maintained

### 3. Priority Levels
- `P0`: Critical path - blocks other development
- `P1`: High priority - needed for MVP
- `P2`: Important but not blocking
- `P3`: Nice to have
- `P4`: Low priority

### 4. Implementation Requirements
1. Every stub must panic with a descriptive message
2. Include input validation
3. Document expected behavior
4. Include example usage

### Example Stub

```go
// STUB: Core - ExecuteModule
// Status: Not Started
// Priority: P0
// Dependencies: ModuleLoader, SessionManager
// Blockers: None
// Last Updated: 2025-07-07
//
// Description: Executes a loaded module with the given options
//
// TODO: Implement module execution logic
// NOTE: Must handle both synchronous and asynchronous execution
// WARNING: Not thread-safe in current implementation
func (m *Module) Execute(options map[string]interface{}) (Result, error) {
    panic("STUB: Module.Execute not implemented")
    
    // Expected behavior:
    // 1. Validate module is loaded and runnable
    // 2. Apply options and validate parameters
    // 3. Execute module logic
    // 4. Return results or error
    
    // Example usage:
    // result, err := module.Execute(map[string]interface{}{
    //     "RHOSTS": "192.168.1.1",
    //     "RPORT": 445,
    // })
}
```

### 5. Tracking Stubs
All stubs must be tracked in `docs/PROGRESS.md` under the appropriate section.

### 6. Removing Stubs
When implementing a stubbed function:
1. Remove the stub header
2. Update the progress tracker
3. Add proper documentation
4. Include tests
5. Update any dependent code
