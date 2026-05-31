---
title: "Constant & Magic Value Management"
domains: [constants, magic-values, paths, timeouts, colors, permissions, sentinel-values]
---

# Constant & Magic Value Management

## Overview

This document defines the classification, extraction, and centralization rules for all non-enum constant values in the forge-cli codebase. Enum constants (status, priority, surface type, etc.) are governed by [enum-constants.md](enum-constants.md).

A **magic value** is any literal (string, number, duration, octal) embedded directly in production logic rather than referenced through a named constant. Magic values harm readability, maintainability, and refactor safety: changing a path or timeout requires auditing every file instead of updating one declaration.

## Classification of Constants

Constants fall into five categories. Each has distinct extraction criteria and centralization rules.

### 1. Path Constants

**Definition**: Filesystem paths and path fragments referenced in production code. Includes directory names, file names, and composite paths built via `filepath.Join`.

**Extraction threshold**: Any path string that appears more than once OR that represents a project-internal convention (not an arbitrary user-supplied value) must be a named constant.

**Examples from codebase**:

| Magic Value | Location | Category |
|---|---|---|
| `"tests/results/raw-output.txt"` | `quality_gate.go:255,284` | Test output path |
| `"tests/results/unit-raw-output.txt"` | `quality_gate.go:159,186,507` | Test output path |
| `"tests/results/"` | `init.go:46` (gitignore entry) | Gitignore path |
| `"/health"` | `serverprobe.go:31` | Default health-check path |

**Target state**: All such paths defined as `const` in the relevant package's `constants.go` or in the existing `pkg/feature/constants.go` for shared paths.

**Rationale**: `pkg/feature/constants.go` already centralizes the majority of path constants (`TestBaseDir`, `TestResultsDir`, `ForgeDir`, etc.). The test output file names (`raw-output.txt`, `unit-raw-output.txt`) are the gap -- they appear as inline string literals in `quality_gate.go` despite being project conventions.

### 2. Color Constants

**Definition**: Hex color codes, ANSI escape sequences, and lipgloss color names used for terminal output styling.

**Extraction threshold**: Any color value used for UI display must be a named constant or style variable in a single location per package.

**Examples from codebase**:

| Magic Value | Location | Category |
|---|---|---|
| `"#7DCFFF"` | `init.go:217` (modeHighlight) | Hex color |
| `"#7DCFFF"` | `init_surfaces.go:20` (surfaceStyle) | Hex color (duplicated) |
| `"#FF8700"` | `init_surfaces.go:17` (conflictStyle) | Hex color |
| `"#9ECE6A"` | `init_surfaces.go:23` (sourceStyle) | Hex color |
| `"\033[33m"` / `"\033[0m"` | `list.go:181` | ANSI escape (yellow) |
| `"green"`, `"yellow"`, `"red"`, `"gray"` | `tree.go:217-231` (statusColor) | Named colors |

**Target state**: All hex color strings and ANSI codes extracted to package-level style constants (e.g., in `internal/cmd/styles.go` or `internal/cmd/constants.go`). The `statusColor` function in `tree.go` is already centralized and acceptable -- named color strings inside a single mapping function are not magic values.

**Rationale**: The duplicate `"#7DCFFF"` across `init.go` and `init_surfaces.go` demonstrates the risk: if the design language changes, two files must be updated independently. Centralizing ensures a single source of truth.

### 3. Timeout and Duration Constants

**Definition**: `time.Duration` literals used for retry delays, probe timeouts, lock acquisition windows, and backoff intervals.

**Extraction threshold**: Any `time.Duration` expression in production code must be a named `const` unless it is a one-off value in test-only code.

**Examples from codebase**:

| Magic Value | Location | Category |
|---|---|---|
| `5 * time.Second` | `quality_gate.go:351` (probe retry interval) | Probe interval |
| `5 * time.Second` | `serverprobe.go:61` (probe timeout) | Probe timeout (duplicated value) |
| `5 * time.Second` | `lock.go:16` (lock timeout) | Lock acquisition timeout |
| `50 * time.Millisecond` | `lock.go:55` (backoff interval) | Lock retry backoff |

**Target state**: Each distinct semantic timeout defined as a named `const` in the owning package. When the same duration value serves different purposes (probe timeout vs lock timeout vs retry interval), each must have its own constant with a descriptive name -- do not share a constant merely because the numeric value coincides.

**Rationale**: The value `5 * time.Second` appears in three packages with three different semantics. Merging them into one shared constant would couple unrelated concerns. Each package defines its own constant with a name that reflects its purpose.

### 4. Sentinel Values

**Definition**: Numeric constants used to represent boundary conditions such as "unreachable", "maximum", "default sort depth", or "no result".

**Extraction threshold**: Any integer literal used as a sentinel (not a count, index, or math operand) must be a named constant with a comment explaining its semantics.

**Examples from codebase**:

| Magic Value | Location | Category |
|---|---|---|
| `99999` | `list.go:442` (fallback sort key) | Sort fallback |
| `99999` | `claim.go:376` (cycle task depth) | Unreachable depth |

**Target state**: Each sentinel extracted to a named `const` with a descriptive name and doc comment.

```go
// fallbackSortPriority is assigned to IDs that cannot be parsed,
// ensuring they sort after all valid business IDs.
const fallbackSortPriority = 99999

// unreachableDepth is assigned to tasks in dependency cycles,
// indicating they are not reachable from any root.
const unreachableDepth = 99999
```

### 5. Permission Constants

**Definition**: Octal permission mode values passed to `os.MkdirAll`, `os.WriteFile`, `os.OpenFile`, etc.

**Extraction threshold**: Permission values are **not extracted** to named constants. Go convention uses inline octal literals (`0o755`, `0o644`) which are universally understood by Go developers. Extracting them to constants like `dirPerm` or `filePerm` adds indirection without clarity.

**Standard values**:
- `0o755` -- directories (owner rwx, group/other rx)
- `0o644` -- files (owner rw, group/other r)

**Rationale**: These values are idiomatic Go. The linter and code review conventions already ensure consistency. No deviation exists.

## Extraction Rules

### When to Extract

Extract a magic value to a named constant when **any** of the following conditions hold:

1. **Duplication**: The same value appears in two or more locations in production code (excluding test files).
2. **Semantic coupling**: The value represents a project convention or domain rule (e.g., a file path pattern, a default timeout) rather than an arbitrary local value.
3. **Tuning surface**: The value is likely to change during tuning or configuration (e.g., retry counts, timeout durations).
4. **Clarity**: The value's purpose is not immediately obvious from context (e.g., `99999` vs `maxRetries`).

### When NOT to Extract

Do **not** extract a value when:

1. The value is a Go idiom (e.g., `0o755` for directory permissions, `os.O_RDWR|os.O_CREATE` for file open flags).
2. The value appears only once and its meaning is clear from surrounding code (e.g., `2` in a context window calculation).
3. The value is in test code (tests may use literals for readability and to avoid coupling to production constants).

### Centralized Management Location

| Category | Location | Rationale |
|---|---|---|
| Shared path constants | `pkg/feature/constants.go` | Already the canonical location for project path conventions. `TestBaseDir`, `TestResultsDir`, `ForgeDir` etc. already live here. |
| Package-local paths | `<package>/constants.go` | Paths used within a single package (e.g., test output file names in `pkg/testrunner`) stay local. |
| Color and style constants | `internal/cmd/styles.go` (new) or `internal/cmd/constants.go` | All CLI display colors centralized in the command package where they are used. |
| Timeout constants | `<package>/constants.go` | Each package owns its timeout semantics. No cross-package sharing of timeout constants. |
| Sentinel values | `<package>/constants.go` | Sentinel semantics are domain-specific to the package. |
| Retry/tuning parameters | `<package>/constants.go` | Retry counts and intervals belong to the package that implements the retry logic. |

**Naming convention for constants files**: `constants.go` at the package root. When a package has both enum-like types and non-enum constants, enum types go in their own files (e.g., `status.go`) while non-enum constants go in `constants.go`.

## Target State Definition

The target state is a codebase where:

1. Every production `.go` file contains zero magic strings for paths, zero inline duration literals for timeouts, and zero unexplained numeric literals for sentinel values.
2. `pkg/feature/constants.go` holds all shared path constants (current + test output file names).
3. Each package with timeout/retry logic has a `constants.go` (or inline `const` block) with named duration constants.
4. CLI color values are centralized in one file within `internal/cmd/`.
5. Sentinel values have descriptive names and doc comments.

## Deviation Analysis

The following table catalogs all magic values identified in the Evidence sources (Reference Files) and additional audit, with their current state and required remediation.

### Path Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| P1 | `"tests/results/raw-output.txt"` | `quality_gate.go:255,284` | Inline string, duplicated 2x | Extract to `pkg/feature/constants.go` or `pkg/testrunner/constants.go` |
| P2 | `"tests/results/unit-raw-output.txt"` | `quality_gate.go:159,186,507` | Inline string, duplicated 3x | Extract to `pkg/feature/constants.go` or `pkg/testrunner/constants.go` |
| P3 | `"tests/results/"` | `init.go:46` | Inline in gitignore entries list | Acceptable -- part of a static gitignore entries list, not a runtime path. No extraction needed. |
| P4 | `"/health"` | `serverprobe.go:31` | Inline default path | Extract to named constant `defaultHealthPath` in `pkg/serverprobe/` |

### Color Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| C1 | `"#7DCFFF"` | `init.go:217` | Inline in lipgloss style | Centralize to `internal/cmd/styles.go` as `colorModeHighlight` |
| C2 | `"#7DCFFF"` | `init_surfaces.go:20` | Inline, duplicated from C1 | Reference shared constant from `internal/cmd/styles.go` |
| C3 | `"#FF8700"` | `init_surfaces.go:17` | Inline in lipgloss style | Centralize to `internal/cmd/styles.go` as `colorConflict` |
| C4 | `"#9ECE6A"` | `init_surfaces.go:23` | Inline in lipgloss style | Centralize to `internal/cmd/styles.go` as `colorSource` |
| C5 | `"\033[33m"` / `"\033[0m"` | `list.go:181` | Raw ANSI escape codes | Replace with lipgloss style or named ANSI constants |

### Timeout Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| T1 | `5 * time.Second` | `quality_gate.go:351` | Inline in `probeWithRetry` call | Extract to `const probeRetryInterval` in `internal/cmd/` |
| T2 | `5 * time.Second` | `serverprobe.go:61` | Inline in `ProbeServers` | Extract to `const defaultProbeTimeout` in `pkg/serverprobe/` |
| T3 | `5 * time.Second` | `lock.go:16` | Already named `defaultLockTimeout` | **No deviation** -- follows convention. |
| T4 | `50 * time.Millisecond` | `lock.go:55` | Inline in `time.Sleep` | Extract to `const lockRetryBackoff` in `pkg/index/` |

### Sentinel Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| S1 | `99999` | `list.go:442` | Inline in `sortKey` | Extract to `const fallbackSortPriority` with doc comment |
| S2 | `99999` | `claim.go:376` | Inline in cycle depth assignment | Extract to `const unreachableDepth` with doc comment |

### Retry Parameter Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| R1 | `3` (retry count) | `quality_gate.go:351` | Inline in `probeWithRetry` call | Extract to `const maxProbeRetries` in `internal/cmd/` |
| R2 | `3` (`maxFixTasksPerStep`) | `quality_gate.go:34` | Already named constant | **No deviation** -- follows convention. |
| R3 | `5` (concise error lines) | `quality_gate.go:169,187` | Inline in `ExtractConciseError` calls | Extract to `const conciseErrorMaxLines` in `internal/cmd/` |
| R4 | `10` (max source files) | `quality_gate.go:603` | Inline in extractSourceFiles | Extract to `const maxSourceFiles` in `internal/cmd/` |

### Permission Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| PM1 | `0o755` / `0755` | Multiple files | Inline octal (Go idiom) | **No deviation** -- inline octal is idiomatic Go. |
| PM2 | `0o644` / `0644` | Multiple files | Inline octal (Go idiom) | **No deviation** -- inline octal is idiomatic Go. |

**Note on mixed octal formats**: `0o755` (Go 1.13+ explicit octal) and `0755` (legacy C-style) are both present. The codebase should standardize on `0o755` format for consistency, but this is a style issue, not a constant extraction concern.

## Relationship to enum-constants.md

This document governs **non-enum** constants. Enum-like values (status, priority, surface type, task type) are governed by [enum-constants.md](enum-constants.md). The boundary is:

- **enum-constants.md**: Values that form a closed set of mutually exclusive options, typed as `type X string` or `type X int`.
- **constants.md (this document)**: All other magic values -- paths, durations, colors, sentinels, numeric tuning parameters.

When a value could fit either document (e.g., a string constant that is not enum-like but is also not a path), prefer this document unless the value has the characteristics of an enum (closed set, type-safe comparison, terminal state detection).
