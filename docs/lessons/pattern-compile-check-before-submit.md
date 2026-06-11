---
created: "2026-05-14"
tags: [testing, dependencies]
---

# Task Executor Must Verify Compilation Before Submitting

## Problem

Task executor subagent completed task 2 (info commands) and marked it as `completed`, but the code had two compilation errors:

1. `root.go:42` — `undefined: lessonCmd` (variable referenced but never defined in the file that declares it)
2. `forensic.go:892` — `truncate redeclared in this block` (duplicate function definition)

The dispatcher only discovered the errors via LSP diagnostics after the subagent returned.

## Root Cause

Causal chain (4 levels):

1. **Symptom**: Build breaks after task marked complete; LSP shows red squiggles
2. **Direct cause**: Subagent wrote code that doesn't compile (`lessonCmd` variable not declared, `truncate` function duplicated)
3. **Root cause**: Task executor's quality gate runs `just test` which uses cached test results — it didn't run `go build ./...` or the build cache masked the error
4. **Trigger condition**: Any task that creates new command files and registers them in `root.go` but the subagent forgets to define the command variable

## Solution

The task executor's quality gate must include an explicit `go build ./...` step. Even if `just test` passes (via cache), a compilation check catches:

- Undefined references across packages
- Duplicate declarations
- Missing imports

## Reusable Pattern

**Before `forge task submit`, always run `go build ./...` (or equivalent compile check for the language).** This is the cheapest, fastest safety net.

Quality gate order for Go projects:
1. `go build ./...` — catches undefined references, duplicate declarations
2. `go vet ./...` — catches common mistakes
3. `go test ./...` — catches logic errors

Skipping step 1 means steps 2 and 3 may not even compile.

## Example

```bash
# Task executor should always do:
go build ./...   # MUST pass before proceeding
go test ./...    # then run tests
```

## Related Files

- `forge-cli/internal/cmd/root.go` — command registration
- `forge-cli/internal/cmd/lesson.go` — was missing, caused `undefined: lessonCmd`
