---
id: "5"
title: "Move forge-config schema and example YAML to CLI, update test paths"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "refactor"
mainSession: false
---

# 5: Move forge-config schema and example YAML to CLI, update test paths

## Description
Move `forge-config.schema.json` and `forge-config.example.yaml` from `plugins/forge/references/shared/` to the `forge-cli/` directory (alongside the test that consumes them), and update the hardcoded relative paths in `config_schema_test.go`.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `forge-config.schema.json` and `forge-config.example.yaml` exist in their new location under `forge-cli/`
- [ ] `config_schema_test.go` paths updated to read from new location
- [ ] `go test ./forge-cli/internal/cmd/ -run TestConfigSchema -v` passes
- [ ] Old files removed from `plugins/forge/references/shared/`

## Hard Rules
- The Go test file uses relative paths from `forge-cli/internal/cmd/` — update the `filepath.Join` calls to point to the new location
- Run `go test` to verify the move before considering the task complete

## Implementation Notes
- Current path resolution: `filepath.Join("..", "..", "..", "plugins", "forge", "references", "shared", "forge-config.schema.json")` from `forge-cli/internal/cmd/`
- Suggested new location: `forge-cli/internal/cmd/testdata/` or `forge-cli/testdata/` — conventional Go test data location
- After moving, update the test helper `schemaPath()` and the `examplePath` in line 197
- Also update `scripts/version.txt` per CLI CLAUDE.md rules (patch bump for dead code path cleanup)
