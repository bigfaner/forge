---
created: "2026-06-10"
tags: [error-handling, dependencies]
---

# `forge justfile scaffold --aggregate` silently fails without `--type`

## Problem

`forge justfile scaffold --aggregate` exits with code 1 and produces zero output — no stdout, no stderr, no error message. Appears completely broken.

## Root Cause

Three conditions combine to produce a silent failure:

1. **`MarkFlagRequired("type")` is unconditional** (`register.go:40`) — Cobra rejects the command before `RunE` executes, because `--type` is missing
2. **`SilenceErrors: true` + `SilenceUsage: true`** (`register.go:27-28`) — Cobra suppresses both the error message and usage hint that would normally explain the missing flag
3. **`rootCmd.Execute()` → `os.Exit(1)`** (`root.go:65`) — any error from Cobra becomes a silent exit

The causal chain: `--aggregate` mode does not need `--type` (aggregate reads surfaces from config), but the flag requirement is enforced at the Cobra level before the command logic can branch on `scaffoldAggregate`.

## Solution

Workaround: pass `--type cli` (any valid type) alongside `--aggregate`:

```bash
forge justfile scaffold --type cli --aggregate
```

Proper fix: remove `Cmd.MarkFlagRequired("type")` from `init()`, move validation into `runScaffold` on the non-aggregate path only (where `ValidateArgs()` already checks it at line 54).

## Reusable Pattern

When a cobra command has mutually exclusive modes (e.g. `--type` vs `--aggregate`), do not use `MarkFlagRequired` unconditionally. Either:
- Validate required flags manually in `RunE` after branching on the mode
- Use cobra's `RegisterFlagCompletionFunc` or `MarkFlagRequired` in a `PreRunE` that checks which mode is active

Additionally, `SilenceErrors: true` without a custom error handler creates silent failures. Pair it with explicit error output or only use it when the command handles all errors internally.

## Related Files

- `forge-cli/internal/cmd/scaffold/register.go` — flag registration, `runScaffold`, `runAggregate`
- `forge-cli/internal/cmd/root.go` — `Execute()` with `os.Exit(1)`
- `forge-cli/internal/cmd/scaffold/generate.go` — `GenerateAggregate()` (works correctly)
