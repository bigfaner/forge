---
name: scope-assignment
description: Path classification rules for task scope assignment (frontend/backend/all)
---

# Scope Assignment Rules

## Path Classification

| Path Pattern | Scope |
|---|---|
| `ui/`, `components/`, `pages/`, `styles/`, `public/` | `frontend` |
| Directory containing only `package.json` | `frontend` |
| `cmd/`, `internal/`, `pkg/`, `api/` | `backend` |
| Directory containing `go.mod`, `Cargo.toml`, or `pyproject.toml` | `backend` |

## Special Cases

- `src/` directory: check for language markers — `go.mod`/`Cargo.toml` without `package.json` → backend; reverse → frontend; both or neither → `undetermined`
- Any path not matching above patterns → `undetermined`

## Computation

1. Classify each affected file path independently
2. All frontend → `"frontend"`
3. All backend → `"backend"`
4. Mixed or any undetermined → `"all"`
5. Non-mixed projects (single-language) always use `"all"`
