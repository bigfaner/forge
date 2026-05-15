---
id: "1"
title: "Remove run.* and graduate.* fields from all 6 manifest.yaml files"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Remove run.* and graduate.* fields from all 6 manifest.yaml files

## Description

All 6 profile manifest.yaml files define `run` and `graduate` command fields that are never parsed by any Go code. These fields duplicate the commands already canonically defined in `justfile-recipes`. Remove them to eliminate the drift source and establish justfile-recipes as the single source of truth.

## Reference Files
- `docs/proposals/justfile-canonical-e2e/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/profile/profiles/go-test/manifest.yaml` | Remove `run` and `graduate` sections |
| `forge-cli/pkg/profile/profiles/java-junit/manifest.yaml` | Remove `run` and `graduate` sections |
| `forge-cli/pkg/profile/profiles/maestro/manifest.yaml` | Remove `run` and `graduate` sections |
| `forge-cli/pkg/profile/profiles/pytest/manifest.yaml` | Remove `run` and `graduate` sections |
| `forge-cli/pkg/profile/profiles/rust-test/manifest.yaml` | Remove `run` and `graduate` sections |
| `forge-cli/pkg/profile/profiles/web-playwright/manifest.yaml` | Remove `run` and `graduate` sections |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] All 6 manifest.yaml files contain zero `run.*` fields (`run.command`, `run.compile`, `run.result-format`)
- [ ] All 6 manifest.yaml files contain zero `graduate.*` fields (`graduate.target-directory`, `graduate.merge-strategy`, `graduate.import-rewrite`, `graduate.compile-check`, `graduate.list-tests`)
- [ ] Remaining manifest fields (name, display, language, file-extension, test-directory, capabilities, templates) are untouched
- [ ] `go test ./...` passes after removal
- [ ] `forge profile get <profile> --manifest` still works for all 6 profiles

## Implementation Notes

- The `graduate.target-directory` field value is always `tests/e2e/` across all profiles. This path is already hardcoded in the graduate strategy .md files and in `pkg/e2e/actions.go` where needed. No code reads the manifest value, so removal is safe.
- Verify by searching the codebase for any code that reads these fields: `grep -r "run\\.command\\|run\\.compile\\|run\\.result-format\\|graduate\\." forge-cli/pkg/` should return zero hits outside the YAML files themselves.
