---
id: "7"
title: "Fix breakdown-tasks scope heuristic for backend projects"
priority: "P2"
estimated_time: "30m"
dependencies: []
type: "enhancement"
scope: "all"
breaking: false
mainSession: false
---

# 7: Fix breakdown-tasks scope heuristic for backend projects

## Description

`breakdown-tasks/SKILL.md:290` classifies `src/` as `frontend` unconditionally. This is wrong for Go/Rust backend projects where `src/` is the standard source directory. The current heuristic has a partial fix (checks for `go.mod`/`Cargo.toml`), but it only applies when `package.json` is also present — pure backend projects without `package.json` still get misclassified.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W4, item 12)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Line 290

## Affected Files

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Fix scope heuristic at line 290 to check for Go/Rust indicators before classifying `src/` as frontend |

## Acceptance Criteria

- [ ] `src/` is classified as `backend` when `go.mod` or `Cargo.toml` exists at the same level (without requiring `package.json` to be absent)
- [ ] `src/` remains `frontend` for Node.js/React projects (when `package.json` exists without `go.mod`/`Cargo.toml`)
- [ ] The heuristic handles mixed projects (both `package.json` and `go.mod`) correctly → `scope: "all"`

## Hard Rules

- Only modify the scope heuristic section (around line 290). Do not change other sections.
- Keep the existing classification algorithm structure (classify → compute scope → write).
- The fix should be a condition addition, not a rewrite of the heuristic.

## Implementation Notes

- Current heuristic (line 290): `frontend: path starts with ... src/ ... or any directory containing package.json with no go.mod/Cargo.toml at the same level`
- Problem: `src/` is listed unconditionally in the frontend path prefixes. The `go.mod`/`Cargo.toml` check only applies to the `package.json` fallback clause.
- Fix: Move `src/` out of the unconditional frontend list. Add a conditional check: if `go.mod` or `Cargo.toml` exists at project root → classify `src/` as `backend`; if `package.json` exists → classify `src/` as `frontend`; otherwise → `undetermined`.
