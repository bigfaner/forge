---
title: "Enum & Constant Organization"
domains: [constants, enums, types, status, surface-type, priority, magic-values, paths, timeouts, colors, sentinel-values, permissions]
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

---

## Non-Enum Constant Management

The following rules govern constants that are **not enum-like** (do not form a closed set of mutually exclusive options) but still require centralized management. Full classification rules, extraction criteria, and deviation analysis are in [constants.md](constants.md).

### TECH-const-001: Path Constants Must Be Named

**Requirement**: Any filesystem path string that appears more than once in production code, or that represents a project-internal convention, must be defined as a named `const` in the appropriate `constants.go` file.

**Shared paths** (referenced by multiple packages) go in `pkg/feature/constants.go`:
```go
// pkg/feature/constants.go
const (
    TestOutputFileName       = "raw-output.txt"
    UnitTestOutputFileName   = "unit-raw-output.txt"
)
```

**Package-local paths** (used within one package) go in `<package>/constants.go`:
```go
// pkg/serverprobe/constants.go
const defaultHealthPath = "/health"
```

**Deviation** (from Evidence -- `quality_gate.go`):
- `"tests/results/raw-output.txt"` appears as inline literal at lines 255 and 284
- `"tests/results/unit-raw-output.txt"` appears as inline literal at lines 159, 186, and 507
- Both should reference constants from `pkg/feature/constants.go` (or `pkg/testrunner/constants.go`)

**Scope**: [CROSS]

### TECH-const-002: Timeout and Duration Values Must Be Named

**Requirement**: Every `time.Duration` expression in production code must be a named `const` with a descriptive name reflecting its semantic purpose (e.g., `probeRetryInterval`, not `timeout5s`). When the same numeric duration serves different purposes, each must have its own constant.

**Pattern**:
```go
// pkg/serverprobe/constants.go
const defaultProbeTimeout = 5 * time.Second

// internal/cmd/constants.go
const probeRetryInterval = 5 * time.Second
const maxProbeRetries    = 3
```

**Anti-pattern** -- do NOT share a constant between unrelated concerns:
```go
// BAD: lock timeout and probe timeout share a constant because both are 5s
const fiveSeconds = 5 * time.Second  // semantically distinct!
```

**Deviation** (from Evidence -- `quality_gate.go`):
- `5 * time.Second` passed inline to `probeWithRetry()` at line 351 -- should be `const probeRetryInterval`
- Retry count `3` passed inline to `probeWithRetry()` at line 351 -- should be `const maxProbeRetries`
- `5 * time.Second` hardcoded in `serverprobe.go:61` -- should be `const defaultProbeTimeout`
- `50 * time.Millisecond` hardcoded in `lock.go:55` -- should be `const lockRetryBackoff`

**Note**: `defaultLockTimeout` in `lock.go:16` already follows this convention -- no deviation.

**Scope**: [CROSS]

### TECH-const-003: Color Values Must Be Centralized

**Requirement**: All hex color codes (`#RRGGBB`), ANSI escape sequences, and lipgloss color names used for terminal display must be centralized in a single file within `internal/cmd/`. This ensures the design language can be updated in one place.

**Pattern**:
```go
// internal/cmd/styles.go
const (
    colorModeHighlight = "#7DCFFF"
    colorConflict      = "#FF8700"
    colorSource        = "#9ECE6A"
    colorCycleMarker   = "\033[33m"
    colorReset         = "\033[0m"
)
```

**Exception**: Named color strings inside a single mapping function (e.g., `statusColor()` returning `"green"`, `"red"`, `"yellow"`, `"gray"`) are acceptable without extraction. The function itself is the centralization point.

**Deviation** (from Evidence -- `init.go`, `init_surfaces.go`, `list.go`):
- `"#7DCFFF"` duplicated across `init.go:217` and `init_surfaces.go:20`
- `"#FF8700"` hardcoded in `init_surfaces.go:17`
- `"#9ECE6A"` hardcoded in `init_surfaces.go:23`
- Raw ANSI codes `"\033[33m"` and `"\033[0m"` in `list.go:181`

**Scope**: [CROSS]

### TECH-const-004: Sentinel Values Must Be Named and Documented

**Requirement**: Any numeric literal used as a sentinel value (representing "unreachable", "maximum", "no result", etc.) must be extracted to a named `const` with a doc comment explaining its semantics.

**Pattern**:
```go
// fallbackSortPriority is assigned to task IDs that cannot be parsed,
// ensuring they sort after all valid business IDs.
const fallbackSortPriority = 99999

// unreachableDepth is assigned to tasks in dependency cycles,
// indicating they are not reachable from any root in BFS traversal.
const unreachableDepth = 99999
```

**Deviation** (from Evidence -- `list.go`, `claim.go`):
- `99999` hardcoded in `list.go:442` as fallback sort priority
- `99999` hardcoded in `claim.go:376` as cycle task depth

**Scope**: [CROSS]

### TECH-const-005: Permission Values Stay Inline

**Requirement**: Octal permission values (`0o755`, `0o644`) are NOT extracted to named constants. Inline octal literals are idiomatic Go and universally understood. Standardize on the `0o` prefix format (Go 1.13+) for consistency, but do not wrap in named constants.

**No deviation** -- current usage follows this rule. Mixed `0o755` / `0755` formatting is a style concern, not a constant management concern.

**Scope**: [CROSS]
