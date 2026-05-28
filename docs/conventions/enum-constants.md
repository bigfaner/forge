---
title: "Enum & Constant Organization"
domains: [constants, enums, types, status, surface-type, priority, magic-values]
---

# Enum & Constant Organization

### TECH-enum-001: Use Typed Constants for All Enum-Like Values

**Requirement**: All enum-like string values (Status, SurfaceType, Priority, etc.) must be defined as typed constants, not used as raw string literals. A typed constant provides compile-time type safety: the compiler rejects passing an arbitrary `string` where a `types.Status` is expected.

**Scope**: [CROSS]

**Pattern to avoid**:
```go
// No type safety — any string accepted
func (t *Task) SetStatus(status string) error { ... }

// Magic value — typo-prone
t.Status = "completed"
if t.Status == "complteed" { ... } // compiles, fails silently
```

**Required pattern**:
```go
// pkg/types/status.go
type Status string

const (
    StatusPending    Status = "pending"
    StatusCompleted  Status = "completed"
    // ...
)

// Usage — type-safe, typo caught at compile time
func (t *Task) SetStatus(status types.Status) error { ... }
t.Status = types.StatusCompleted
```

### TECH-enum-002: Centralize Enums in `pkg/types/` (Leaf Package)

**Requirement**: All shared enum types and constants live in `pkg/types/`. This package must not import any other forge-cli internal package — it is a pure type definition module. Each enum category gets its own file: `status.go`, `surface.go`, `priority.go`.

**Rationale**: `pkg/types/` is imported by `pkg/feature`, `pkg/task`, `pkg/forgeconfig`, and `internal/cmd`. As a leaf package, it breaks no dependency cycles.

**Scope**: [CROSS]

### TECH-enum-003: Provide Enumeration Helpers

**Requirement**: Each enum type must provide two standard helpers:

- `AllXxx() []Xxx` — returns all valid values (replaces ad-hoc `map[string]bool` validation maps)
- `IsTerminalXxx(x Xxx) bool` — if the enum has terminal states (currently only `Status`)

**Pattern**:
```go
func AllStatuses() []Status {
    return []Status{StatusPending, StatusInProgress, StatusCompleted, ...}
}

func IsTerminalStatus(s Status) bool {
    return s == StatusCompleted || s == StatusSkipped || s == StatusRejected
}
```

**Scope**: [CROSS]

### TECH-enum-004: Convert at Boundaries, Not Internally

**Requirement**: External interfaces (CLI flags, config parsing, JSON/YAML) use `string`. Convert to the typed constant at the boundary entry point — internal code never sees raw strings.

**Pattern**:
```go
// Boundary: CLI flag → typed constant
status := types.Status(cobraFlag)

// Boundary: config → typed constant
surface := types.SurfaceType(viper.GetString("surfaceType"))

// Internal: always typed
func processTask(status types.Status) { ... }
```

**Note**: Go's `type X string` preserves JSON/YAML marshal/unmarshal behavior — no custom marshalers needed.

**Scope**: [CROSS]

### TECH-enum-005: Re-export for Backward Compatibility

**Requirement**: When migrating constants from their original package (e.g., `pkg/feature/constants.go`), provide type aliases or variable re-exports to avoid breaking downstream code during transition.

**Pattern**:
```go
// pkg/feature/constants.go — re-export
type Status = types.Status

var (
    StatusPending    = types.StatusPending
    StatusCompleted  = types.StatusCompleted
    // ...
)
```

**Scope**: [CROSS]

### TECH-enum-006: Zero Magic Values in Production Code

**Requirement**: Production `.go` files must contain zero raw string literals matching enum values. All references use the typed constants defined in `pkg/types/`. Test files may use string literals only when testing serialization or boundary conversion.

**Scope**: [CROSS]

### TECH-enum-007: Enums Not to Migrate

**Requirement**: Some enum-like constants are tightly coupled to their domain logic and should NOT be moved to `pkg/types/`:

- **Task Type** constants (`TypeCodingFeature`, etc.) — remain in `pkg/task/types.go` due to deep coupling with task logic
- **Config dotpath** keys (`"eval.proposal"`, etc.) — not enum values, but nested config paths
- **Path** constants (`"prd"`, `"design"`, etc.) — not enums, belong to their own domain

**Scope**: [CROSS]
