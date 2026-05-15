---
id: "2"
title: "Delegate e2e/actions.go functions to just recipes"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 2: Delegate e2e/actions.go functions to just recipes

## Description

Replace the profile-specific switch/case dispatch in `e2e/actions.go` with thin `exec.Command("just", ...)` delegation. Run, Setup, Compile, and Discover each become a single `exec.Command` call that delegates to the corresponding justfile recipe. Verify remains unchanged (generic file scanner).

## Reference Files
- `docs/proposals/justfile-canonical-e2e/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/e2e/actions.go` | Replace profile switch/case in Run/Setup/Compile/Discover with `exec.Command("just", ...)` delegation; add `just` not-found error check |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `Run()` calls `exec.Command("just", "test-e2e")` with optional `--feature <name>` argument
- [ ] `Setup()` calls `exec.Command("just", "e2e-setup")`
- [ ] `Compile()` calls `exec.Command("just", "e2e-compile")`
- [ ] `Discover()` calls `exec.Command("just", "e2e-discover")`
- [ ] `Verify()` is completely unchanged
- [ ] Non-zero exit from `just` produces a non-nil error; zero exit produces nil error
- [ ] Forge CLI exit code matches the `just` exit code
- [ ] When `just` is not on PATH, error message: `error: 'just' is required but not found on PATH. Install: https://github.com/casey/just`
- [ ] All 6 profiles work without per-profile code paths (justfile recipes handle profile dispatch)
- [ ] `go build ./...` passes

## Hard Rules

- Do NOT modify `Verify()` — it is profile-agnostic and has no justfile equivalent.
- The justfile recipes already read the active profile from `.forge/config.yaml` at runtime. The Go code must NOT pass profile information to `just` — the recipe handles it internally.
- `Run()` must pass `opts.Feature` as a justfile argument when non-empty. The format is `just test-e2e feature=<name>` (just named argument syntax) or equivalent positional argument depending on the recipe signature.

## Implementation Notes

- The existing `runner` variable and `ExecRunner` interface in `exec.go` can remain — just change the command from `go test`/`npx playwright` to `just <recipe>`.
- The `formatToolError` helper already formats subprocess errors with the command name and first stderr line. It will work as-is with `just` as the command name.
- Feature filtering: current `go-test` Run() filters by directory `./tests/e2e/features/{Feature}/...`. The justfile recipe `test-e2e` accepts a `feature` parameter and applies equivalent filtering. This is a behavior change from directory-scoped (Go) to test-name-based (justfile) filtering for the `go-test` profile. Document this in the commit message.
- The `ResolveProfile()` function in `e2e.go` is still needed — it validates that a profile is configured before delegating to `just`. This prevents confusing "no justfile found" errors when the real issue is no profile configured.
