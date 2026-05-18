---
id: "5"
title: "Move forge-config schema and example YAML to CLI, update test paths"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 5: Move forge-config schema and example YAML to CLI, update test paths

## Description
Move `forge-config.schema.json` and `forge-config.example.yaml` from `plugins/forge/references/shared/` to the `forge-cli/` directory (alongside the test that consumes them), and update the hardcoded relative paths in `config_schema_test.go`.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `forge-cli/internal/cmd/config_schema_test.go` — Go test file with hardcoded paths

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/config_schema_test.go` | Update `filepath.Join` path expressions to point to new file location |

### Move
| File | From | To |
|------|------|----|
| `forge-config.schema.json` | `plugins/forge/references/shared/` | `forge-cli/internal/cmd/testdata/` (or `forge-cli/testdata/`) |
| `forge-config.example.yaml` | `plugins/forge/references/shared/` | `forge-cli/internal/cmd/testdata/` (or `forge-cli/testdata/`) |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/references/shared/forge-config.schema.json` | Moved to CLI repo |
| `plugins/forge/references/shared/forge-config.example.yaml` | Moved to CLI repo |

## Acceptance Criteria
- [ ] `forge-config.schema.json` and `forge-config.example.yaml` exist under `forge-cli/internal/cmd/testdata/` (or `forge-cli/testdata/`)
- [ ] `config_schema_test.go` `filepath.Join` expressions updated to resolve from new location
- [ ] `go test ./forge-cli/internal/cmd/ -run TestConfigSchema -v` passes
- [ ] Old files removed from `plugins/forge/references/shared/`
- [ ] No path reference to `plugins/forge/references/` remains in `config_schema_test.go`

## Hard Rules
- The Go test file uses relative paths from `forge-cli/internal/cmd/` — update the `filepath.Join` calls to point to the new location
- Run `go test` to verify the move before considering the task complete

## Implementation Notes
- Current path resolution: `filepath.Join("..", "..", "..", "plugins", "forge", "references", "shared", "forge-config.schema.json")` from `forge-cli/internal/cmd/`
- Suggested new location: `forge-cli/internal/cmd/testdata/` or `forge-cli/testdata/` — conventional Go test data location
- After moving, update the test helper `schemaPath()` and the `examplePath` in line 197
- Choose `testdata/` subdirectory following Go convention for test fixtures
