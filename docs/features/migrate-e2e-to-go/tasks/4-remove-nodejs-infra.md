---
id: "4"
title: "Remove Node.js test infrastructure"
priority: "P2"
estimated_time: "15m"
dependencies: ["2", "3"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Remove Node.js test infrastructure

## Description

After all Playwright tests have been successfully converted to Go and verified, remove the entire Node.js test infrastructure from `tests/e2e/`. This eliminates the dual-stack maintenance burden and the Node.js runtime dependency.

Only run this task after tasks 1-3 are complete and ALL converted Go tests pass.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| None | |

### Modify
| File | Changes |
|------|---------|
| None | |

### Delete
| File | Reason |
|------|--------|
| `tests/e2e/package.json` | Node.js dependency manifest |
| `tests/e2e/playwright.config.ts` | Playwright configuration |
| `tests/e2e/tsconfig.json` | TypeScript configuration |
| `tests/e2e/helpers.ts` | Node.js test helpers |
| `tests/e2e/gen-test-scripts/cli.spec.ts` | Converted to Go in task 1 |
| `tests/e2e/justfile-execution/justfile-execution.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/task-cli/typed-task-dispatch.spec.ts` | Converted to Go in task 3 |
| `tests/e2e/scope-resolution/scope-resolution.spec.ts` | Converted to Go in task 3 |
| `tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/justfile-e2e-integration/mixed-template.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/justfile-e2e-integration/cli.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/init-justfile/init-justfile.spec.ts` | Converted to Go in task 3 |
| `tests/e2e/plugin-content/skill-content.spec.ts` | Converted to Go in task 1 |
| `tests/e2e/features/forge-testing-optimization/cli.spec.ts` | Duplicate of gen-test-scripts, merged in task 1 |
| `tests/e2e/features/forge-testing-optimization/playwright.config.ts` | Feature-scoped Playwright config |
| `tests/e2e/node_modules/` (if exists) | Node.js dependencies |

## Acceptance Criteria
- [ ] No `.spec.ts` files remain in `tests/e2e/` (verified with `find tests/e2e -name '*.spec.ts' | wc -l` → 0)
- [ ] No `package.json` in `tests/e2e/`
- [ ] No `node_modules/` in `tests/e2e/`
- [ ] No `playwright.config.ts` or `tsconfig.json` in `tests/e2e/`
- [ ] Empty directories removed (e.g., `tests/e2e/gen-test-scripts/`, `tests/e2e/task-cli/`, etc.)
- [ ] `go test ./tests/e2e/... -v -tags=e2e` still passes (regression check)

## Hard Rules
- MUST verify all Go tests pass BEFORE deleting any Playwright files
- MUST NOT delete `tests/e2e/fixtures/` directory — these are test data, not Node.js infrastructure
- MUST NOT delete `tests/e2e/results/` directory — these are test results
- MUST NOT modify any files in `plugins/forge/skills/gen-test-scripts/templates/` — these are test generation templates, not the test infrastructure being removed

## Implementation Notes

### Pre-deletion verification:
```bash
# Verify all Go e2e tests pass
cd forge-cli && go test ./tests/e2e/... -v -tags=e2e

# Count Playwright test files to confirm scope
find tests/e2e -name '*.spec.ts' | wc -l
```

### Deletion script:
```bash
# Remove Playwright test files
find tests/e2e -name '*.spec.ts' -delete

# Remove Node.js config files
rm tests/e2e/package.json
rm tests/e2e/playwright.config.ts
rm tests/e2e/tsconfig.json
rm tests/e2e/helpers.ts

# Remove feature-scoped Playwright config
rm -rf tests/e2e/features/forge-testing-optimization/

# Remove node_modules if exists
rm -rf tests/e2e/node_modules/

# Remove empty directories
find tests/e2e -type d -empty -delete
```

### Post-deletion verification:
```bash
# Confirm no .spec.ts remains
find tests/e2e -name '*.spec.ts' | wc -l  # Expected: 0

# Confirm Go tests still pass
cd forge-cli && go test ./tests/e2e/... -v -tags=e2e
```
