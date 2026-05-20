---
feature: "e2e-isolated-build-artifact"
status: tasks
mode: quick
---

# Feature (Quick): e2e-isolated-build-artifact

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/e2e-isolated-build-artifact/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Isolate forge-cli/tests/e2e/ to TestMain auto-build | pending | tasks/1-isolate-forge-cli-e2e.md |
| 2 | Isolate justfile-canonical-e2e/ to TestMain auto-build | pending | tasks/2-isolate-justfile-canonical-e2e.md |
| 3 | Fix tests/e2e/ feature tests to use TestMain-built binary | pending | tasks/3-fix-tests-e2e-features.md |
| 4 | Simplify e2e-setup in justfile to optional cache optimization | pending | tasks/4-simplify-e2e-setup.md |
| 5 | Update TEST-isolation-004 scope to cover all test locations | pending | tasks/5-update-testing-isolation-convention.md |
