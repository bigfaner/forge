---
feature: "migrate-e2e-to-go"
status: completed
mode: quick
---

# Feature (Quick): migrate-e2e-to-go

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/migrate-e2e-to-go/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Add file assertion helpers to testkit | completed | tasks/1-testkit-helpers.md |
| 2 | Convert gen-test-scripts and forge-testing-optimization tests | completed | tasks/2-convert-gen-test-scripts.md |
| 3 | Convert justfile-execution tests | completed | tasks/3-convert-justfile-execution.md |
| 4 | Convert task-cli typed-task-dispatch tests | completed | tasks/4-convert-task-cli.md |
| 5 | Convert scope-resolution tests | completed | tasks/5-convert-scope-resolution.md |
| 6 | Convert justfile-e2e-integration tests (forge-justfile + detection-assembly) | completed | tasks/6-convert-justfile-integration-a.md |
| 7 | Convert justfile-e2e-integration tests (mixed-template + cli) | completed | tasks/7-convert-justfile-integration-b.md |
| 8 | Convert init-justfile and plugin-content tests | completed | tasks/8-convert-init-justfile-plugin-content.md |
| 9 | Remove Node.js test infrastructure | completed | tasks/9-remove-nodejs-infrastructure.md |
| T-quick-1 | Generate test cases from proposal | skipped | tasks/T-quick-1.md |
| T-quick-2 | Generate test scripts from test cases | skipped | tasks/T-quick-2.md |
| T-quick-3 | Execute feature e2e tests | skipped | tasks/T-quick-3.md |
| T-quick-4 | Graduate tests to regression suite | skipped | tasks/T-quick-4.md |
| T-quick-5 | Run full regression suite | skipped | tasks/T-quick-5.md |
