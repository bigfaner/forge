---
id: "9"
title: "Remove Node.js test infrastructure"
priority: "P2"
estimated_time: "15m"
dependencies: ["2", "3", "4", "5", "6", "7", "8"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 9: Remove Node.js test infrastructure

## Description

After all Playwright test cases have been converted to Go and verified passing, remove the entire Node.js test infrastructure. This is the final cleanup step that eliminates the Node.js runtime dependency.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|
| `tests/e2e/package.json` | Node.js dependency manifest |
| `tests/e2e/playwright.config.ts` | Playwright configuration |
| `tests/e2e/tsconfig.json` | TypeScript configuration |
| `tests/e2e/helpers.ts` | Shared test helpers |
| `tests/e2e/helpers.js` | Compiled helpers (if present) |
| `tests/e2e/gen-test-scripts/cli.spec.ts` | Converted to Go in task 2 |
| `tests/e2e/justfile-execution/justfile-execution.spec.ts` | Converted to Go in task 3 |
| `tests/e2e/task-cli/typed-task-dispatch.spec.ts` | Converted to Go in task 4 |
| `tests/e2e/scope-resolution/scope-resolution.spec.ts` | Converted to Go in task 5 |
| `tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts` | Converted to Go in task 6 |
| `tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts` | Converted to Go in task 6 |
| `tests/e2e/justfile-e2e-integration/mixed-template.spec.ts` | Converted to Go in task 7 |
| `tests/e2e/justfile-e2e-integration/cli.spec.ts` | Converted to Go in task 7 |
| `tests/e2e/init-justfile/init-justfile.spec.ts` | Converted to Go in task 8 |
| `tests/e2e/plugin-content/skill-content.spec.ts` | Converted to Go in task 8 |
| `tests/e2e/features/forge-testing-optimization/cli.spec.ts` | Merged into task 2 |

## Acceptance Criteria
- [ ] No `.spec.ts` files remain in `tests/e2e/`
- [ ] No `package.json` in `tests/e2e/`
- [ ] No `node_modules/` in `tests/e2e/`
- [ ] No `playwright.config.ts` in `tests/e2e/`
- [ ] `go test ./tests/e2e/... -v -tags=e2e` still passes (no regressions)
- [ ] `go build ./...` passes

## Hard Rules

- ONLY run this task after ALL conversion tasks (2-8) are verified passing
- Verify `go test ./tests/e2e/... -v -tags=e2e` passes BEFORE deleting any files
- Do NOT delete `tests/e2e/` directory itself — it may contain fixtures or other non-Playwright content
- Do NOT delete `tests/e2e/fixtures/` — test fixtures are still used by Go tests
- Run a final `go build ./...` after deletion to confirm nothing references removed files

## Implementation Notes

- This is the final task. Run full Go test suite before and after deletion to confirm zero regressions.
- Check for any CI configuration references to Playwright/Node.js that should also be cleaned up (outside `tests/e2e/`).
- Verify no other files in the repo import or reference the deleted Playwright test infrastructure.
