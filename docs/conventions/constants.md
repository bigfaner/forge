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

**Examples from codebase** (current state -- all extracted):

| Constant Name | Location | Category |
|---|---|---|
| `TestOutputFileName = "raw-output.txt"` | `pkg/feature/constants.go` | Test output file name |
| `UnitTestOutputFileName = "unit-raw-output.txt"` | `pkg/feature/constants.go` | Test output file name |
| `"tests/results/"` | `init.go:35` (gitignore entry, `gitignoreEntries` slice) | Gitignore path (static list, acceptable) |
| `defaultHealthPath = "/health"` | `pkg/serverprobe/constants.go` | Default health-check path |

**Note**: Some constants exist as inline defaults in `pkg/testrunner/` code rather than as named constants, because they are single-use and contextually clear.

**Target state**: All such paths defined as `const` in the relevant package's `constants.go` or in the existing `pkg/feature/constants.go` for shared paths.

### 2. Color Constants

**Definition**: Hex color codes, ANSI escape sequences, and lipgloss color names used for terminal output styling.

**Extraction threshold**: Any color value used for UI display must be a named constant or style variable in a single location per package.

**Examples from codebase** (current state -- all centralized):

| Constant Name | Location | Category |
|---|---|---|
| `colorModeHighlight = "#7DCFFF"` | `internal/cmd/styles.go` | Hex color (mode keywords) |
| `colorConflict = "#FF8700"` | `internal/cmd/styles.go` | Hex color (conflict text) |
| `colorSource = "#9ECE6A"` | `internal/cmd/styles.go` | Hex color (source text) |
| `"green"`, `"yellow"`, `"red"`, `"gray"` | `tree.go` (statusColor function) | Named colors (acceptable) |

**Target state**: All hex color strings and ANSI codes extracted to package-level style constants in `internal/cmd/styles.go`. The `statusColor` function in `tree.go` is already centralized and acceptable -- named color strings inside a single mapping function are not magic values.

### 3. Timeout and Duration Constants

**Definition**: `time.Duration` literals used for retry delays, probe timeouts, lock acquisition windows, and backoff intervals.

**Extraction threshold**: Any `time.Duration` expression in production code must be a named `const` unless it is a one-off value in test-only code.

**Examples from codebase** (current state -- all extracted):

| Constant Name | Location | Category |
|---|---|---|
| `probeRetryInterval = 5 * time.Second` | `internal/cmd/qualitygate/constants.go` | Probe retry interval |
| `maxProbeRetries = 3` | `internal/cmd/qualitygate/constants.go` | Max probe retry count |
| `defaultProbeTimeout = 5 * time.Second` | `pkg/serverprobe/constants.go` | Probe timeout |
| `defaultLockTimeout = 5 * time.Second` | `pkg/index/lock.go` | Lock acquisition timeout |
| `lockRetryBackoff = 50 * time.Millisecond` | `pkg/index/lock.go` | Lock retry backoff |

**Target state**: Each distinct semantic timeout defined as a named `const` in the owning package. When the same duration value serves different purposes (probe timeout vs lock timeout vs retry interval), each must have its own constant with a descriptive name -- do not share a constant merely because the numeric value coincides.

### 4. Sentinel Values

**Definition**: Numeric constants used to represent boundary conditions such as "unreachable", "maximum", "default sort depth", or "no result".

**Extraction threshold**: Any integer literal used as a sentinel (not a count, index, or math operand) must be a named constant with a comment explaining its semantics.

**Examples from codebase** (current state -- all extracted):

| Constant Name | Location | Category |
|---|---|---|
| `fallbackSortPriority = 99999` | `internal/cmd/task/list.go` | Sort fallback for unparseable IDs |
| `unreachableDepth = 99999` | `internal/cmd/task/claim.go` | Unreachable depth for cycle tasks |

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
| P1 | `"raw-output.txt"` | `pkg/feature/constants.go` | **Fixed**: extracted as `TestOutputFileName` | Done |
| P2 | `"unit-raw-output.txt"` | `pkg/feature/constants.go` | **Fixed**: extracted as `UnitTestOutputFileName` | Done |
| P3 | `"tests/results/"` | `init.go:46` | Acceptable -- part of a static gitignore entries list, not a runtime path. No extraction needed. | N/A |
| P4 | `"/health"` | `pkg/serverprobe/constants.go` | **Fixed**: extracted as `defaultHealthPath` | Done |

### Color Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| C1 | `"#7DCFFF"` | `internal/cmd/styles.go` | **Fixed**: centralized as `colorModeHighlight` | Done |
| C2 | `"#FF8700"` | `internal/cmd/styles.go` | **Fixed**: centralized as `colorConflict` | Done |
| C3 | `"#9ECE6A"` | `internal/cmd/styles.go` | **Fixed**: centralized as `colorSource` | Done |
| C4 | `"\033[33m"` / `"\033[0m"` | `list.go` | **Fixed**: ANSI codes cleaned up | Done |

### Timeout Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| T1 | `5 * time.Second` | `internal/cmd/qualitygate/constants.go` | **Fixed**: extracted as `probeRetryInterval` | Done |
| T2 | `5 * time.Second` | `pkg/serverprobe/constants.go` | **Fixed**: extracted as `defaultProbeTimeout` | Done |
| T3 | `5 * time.Second` | `pkg/index/lock.go` | **No deviation**: already named `defaultLockTimeout` | N/A |
| T4 | `50 * time.Millisecond` | `pkg/index/lock.go` | **Fixed**: extracted as `lockRetryBackoff` | Done |

### Sentinel Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| S1 | `99999` | `internal/cmd/task/list.go` | **Fixed**: extracted as `fallbackSortPriority` with doc comment | Done |
| S2 | `99999` | `internal/cmd/task/claim.go` | **Fixed**: extracted as `unreachableDepth` with doc comment | Done |

### Retry Parameter Deviations

| # | Magic Value | File:Line | Current State | Remediation |
|---|---|---|---|---|
| R1 | `3` (retry count) | `internal/cmd/qualitygate/constants.go` | **Fixed**: extracted as `maxProbeRetries` | Done |
| R2 | `3` (`maxFixTasksPerStep`) | `quality_gate.go:34` | **No deviation**: already named constant | N/A |
| R3 | `5` (concise error lines) | `internal/cmd/qualitygate/constants.go` | **Fixed**: extracted as `conciseErrorMaxLines` | Done |
| R4 | `10` (max source files) | `internal/cmd/qualitygate/constants.go` | **Fixed**: extracted as `maxSourceFiles` | Done |

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
