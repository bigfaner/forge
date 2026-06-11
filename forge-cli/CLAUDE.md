# Mindset

Think from `first principles`, rejecting empiricism and path dependency. Stay cautious. Start from raw requirements — if the goal is unclear, stop and discuss; if the goal is clear but the path is suboptimal, suggest a more direct approach.

All responses must include two parts:

**Direct execution** — deliver the requested result.

**Deep interaction** — `critically challenge` the original request:
  - Question whether the motivation diverges from the goal (XY problem)
  - Analyze drawbacks of the current path
  - Suggest more elegant alternatives

## Development Rules

### TDD (mandatory)

RED → GREEN → REFACTOR cycle. Coverage target **80%+**. Table-driven tests. Mock external dependencies.

### Commit Convention

Use the `/git-commit` SKILL

### Dependency Direction

`cmd → internal → pkg`, reverse is strictly forbidden. Modules interact through interfaces.

## Common Commands

```bash
go build ./...                  # Compile check
go vet ./...                    # Static analysis
go test -race -cover ./...      # Tests (with race detection and coverage)
golangci-lint run ./...         # Lint
bash ../claude-code-go/scripts/lint-arch.sh       # Architecture constraint verification
make check-docs                 # Doc freshness check (validate docs match code)
```

## Documentation

### Index

| Document | Description |
|----------|-------------|
| [docs/OVERVIEW.md](docs/OVERVIEW.md) | Feature overview: core functionality, directory structure, data models |
| [docs/WORKFLOW.md](docs/WORKFLOW.md) | Key workflows: task claiming, hooks, validation, workflows |

### Sync Rules

1. **Feature changes must update docs**: when adding/modifying/removing features, sync `docs/OVERVIEW.md`
2. **Workflow changes must update docs**: when modifying core logic or workflows, sync `docs/WORKFLOW.md`
3. **Pre-commit check**: when code changes affect documented areas, confirm docs are synced
4. **Doc freshness test**: run `make check-docs` after changing structs in `pkg/task/types.go`, constants in `pkg/feature/constants.go`, or detection logic in `internal/cmd/all_completed.go`
