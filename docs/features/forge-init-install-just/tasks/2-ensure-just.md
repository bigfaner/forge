---
id: "2"
title: "Implement ensureJust core logic with detection and installation"
priority: "P1"
estimated_time: "3h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Implement ensureJust core logic with detection and installation

## Description

Implement the core `ensureJust()` function that detects whether `just` is installed, checks its version against the minimum (>= 1.40.0), and installs it either via system package manager or embedded binary fallback.

This is the main logic task — it creates a new package `pkg/just/ensure.go` with the full detect → confirm → install flow.

## Reference Files
- `docs/proposals/forge-init-install-just/proposal.md` — Source proposal
- `forge-cli/pkg/just/just.go` — Existing justfile utilities (quality gate)
- `forge-cli/internal/embedded/just/just_binary.go` — Embedded binary accessor (from task 1)
- `forge-cli/internal/cmd/init.go` — initAction struct and step pattern

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/just/ensure.go` | Core ensureJust logic: detect, version check, install dispatch |
| `forge-cli/pkg/just/ensure_test.go` | Unit tests for ensureJust functions |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `DetectJust() (path string, version string, found bool)` locates `just` in PATH and parses its version
- [ ] `ParseJustVersion(output string) (semver, error)` parses `just --version` output (e.g., `just 1.40.0`)
- [ ] `IsMinimumVersion(version string, minimum string) bool` compares versions (minimum = "1.40.0")
- [ ] `EnsureJust(io.Reader, io.Writer) EnsureResult` orchestrates the full flow with user interaction
- [ ] Package manager installation works for: `brew install just` (macOS), `cargo install just` (cross-platform), `scoop install just` / `choco install just` (Windows)
- [ ] Embedded binary fallback extracts to `~/.forge/bin/just` when package manager fails or is unavailable
- [ ] `EnsureResult` struct carries: Status (INSTALLED/SKIPPED/FAILED), Version, Method (brew/cargo/scoop/choco/embedded), Detail
- [ ] User is prompted for confirmation before installation (unless CI/pipe: check `os.Stdin`)
- [ ] Outdated version (< 1.40.0) triggers upgrade prompt with clear warning
- [ ] All unit tests pass with table-driven test patterns

## Hard Rules

- Do NOT modify `internal/cmd/init.go` — this task only creates the `pkg/just/ensure.go` package
- Package dependency direction: `pkg/just/` may import `internal/embedded/just/` for binary access (this is acceptable since embedded is infrastructure, not business logic)
- User confirmation MUST check that stdin is a terminal — skip prompt and abort if piped

## Implementation Notes

- Use `exec.LookPath("just")` for detection
- Use `exec.Command("just", "--version")` for version check
- Use `runtime.GOOS` and `runtime.GOARCH` to determine platform for package manager dispatch
- For embedded fallback: extract binary from `embedded.JustBinary()` to `~/.forge/bin/just`, set executable permission (chmod 0755 on non-Windows)
- PATH update: print instructions for user to add `~/.forge/bin/` to PATH; do NOT auto-modify shell profiles (per proposal out-of-scope)
- Handle edge case: Windows binary needs `.exe` extension
- The `~/.forge/bin/` directory should be created if it doesn't exist (use `os.UserHomeDir()`)
