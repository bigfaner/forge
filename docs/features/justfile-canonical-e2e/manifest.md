---
feature: "justfile-canonical-e2e"
status: tasks
mode: quick
---

# Feature (Quick): justfile-canonical-e2e

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/justfile-canonical-e2e/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Remove run.* and graduate.* fields from all 6 manifest.yaml files | pending | 1-remove-manifest-command-fields.md |
| 2 | Delegate e2e/actions.go functions to just recipes | pending | 2-delegate-actions-to-just.md |
| 3 | Update e2e/actions_test.go for just delegation | pending | 3-update-actions-tests.md |
| 4 | Version bump to 3.10.0 | pending | 4-version-bump.md |
| T-quick-1 | Generate Quick Test Cases (go-test) | pending | quick-test-cases-go-test.md |
| T-quick-2 | Generate Quick Test Scripts (go-test) | pending | quick-gen-scripts-go-test.md |
| T-quick-3 | Run Quick E2E Tests (go-test) | pending | quick-run-tests-go-test.md |
| T-quick-4 | Graduate Quick Test Scripts (go-test) | pending | quick-graduate-go-test.md |
| T-quick-5 | Verify Quick E2E Regression | pending | quick-verify-regression.md |
