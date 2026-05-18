---
name: eval-test-cases
description: Evaluate test cases for downstream executability. Dispatches per-type eval when per-type files exist, falls back to monolithic eval otherwise.
argument-hint: "[--target 900] [--iterations 6]"
---

# Eval Test Cases

## Step 1: Detect Per-Type Files

1. Resolve feature directory: `docs/features/<current-feature>/testing/`
2. Glob for per-type files: `testing/*-test-cases.md` (i.e., `ui-test-cases.md`, `tui-test-cases.md`, `mobile-test-cases.md`, `api-test-cases.md`, `cli-test-cases.md`)

## Step 2: Dispatch

### Per-Type Mode (per-type files exist)

For each per-type file found in Step 1, invoke the eval skill once:

```
Skill(skill="forge:eval", args="--type {type}-test-cases [--target N] [--iterations N]")
```

Where `{type}` is derived from the filename: `ui-test-cases.md` -> `ui-test-cases`, `api-test-cases.md` -> `api-test-cases`, etc.

Each invocation receives a single `{type}-test-cases.md` file path, NOT the entire `testing/` directory.

After all per-type evals complete, produce an aggregated summary:

```
## Eval-Test-Cases Aggregate Report

| Type | File | Score | Target | Pass? |
|------|------|-------|--------|-------|
| UI   | ui-test-cases.md   | X/1000 | 900 | yes/no |
| API  | api-test-cases.md  | X/1000 | 900 | yes/no |
| **Overall** | | | | **pass/fail** |

Overall passes if every type meets its target.
```

### Legacy Fallback (no per-type files)

If no per-type files exist, fall back to monolithic mode:

```
Skill(skill="forge:eval", args="--type test-cases [--target N] [--iterations N]")
```

This evaluates the single `testing/test-cases.md` file using the legacy `test-cases` rubric.
