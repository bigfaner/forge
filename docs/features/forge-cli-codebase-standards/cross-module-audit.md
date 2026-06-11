# Cross-Module Dependency Audit

**Date:** 2026-05-30
**Module Audited:** `forge-cli` (module path: `forge-cli`)
**Auditor:** Automated audit per task 7

## Audit Method

Two-layer analysis covering Go module level and text-level import detection.

### Layer 1: Go Module Level

| Command | Scope | Result |
|---------|-------|--------|
| `go list -m all` (forge-cli/) | List all dependencies of forge-cli | Only external dependencies (charmbracelet, cobra, testify, etc.). No monorepo sibling modules. |
| `go mod graph` (forge-cli/) | Full dependency graph | All edges originate from `forge-cli` to external packages. No reverse dependency from other monorepo modules. |
| `go mod graph` (tests/) | Check tests module graph | Module `forge-tests` depends only on `stretchr/testify`. Zero references to `forge-cli`. |

### Layer 2: Text-Level Import Search

```bash
grep -rn '"forge-cli/internal\|"forge-cli/pkg' --include='*.go' | grep -v '^forge-cli/'
```

**Result:** No matches. Zero cross-module imports of `forge-cli/internal` or `forge-cli/pkg` found outside `forge-cli/` itself.

## Monorepo Module Inventory

| Module | Path | Module Path | Depends on forge-cli? |
|--------|------|-------------|----------------------|
| forge-cli | `forge-cli/` | `forge-cli` | N/A (self) |
| forge-tests | `tests/` | `forge-tests` | No |

## tests/testkit Integration Pattern

The `tests/testkit/` package interacts with forge-cli at the **binary level** via `exec.Command`, not via Go imports. This is the correct integration pattern for e2e tests:

- `ForgeBinary` is built from source at test init time (`go build -o <tmp> ./cmd/forge`)
- All CLI invocation uses `exec.Command(ForgeBinary, args...)`
- No Go-level import of `forge-cli/internal` or `forge-cli/pkg` packages

## Fallback Evaluation

**Result: No cross-module dependencies detected.**

- No Go module in the monorepo imports `forge-cli` as a dependency
- No Go source file outside `forge-cli/` imports `forge-cli/internal` or `forge-cli/pkg`
- The `tests/` module uses binary-level integration only

**Phase 2c implication:** Phase 2c (removing compatibility layers) can be fully executed without risk of breaking other modules. The `check-cross-module-deps` target in `forge-cli/Makefile` provides ongoing CI protection to ensure this state is maintained.

## CI Protection

A `check-cross-module-deps` target has been added to `forge-cli/Makefile`. It searches for `"forge-cli/internal` imports in Go files outside the `forge-cli/` directory and fails the build if any are found.
