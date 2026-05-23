---
feature: "forge-cli-clean-code"
created: "2026-05-24"
status: tasks
mode: quick
---

# Feature (Quick): forge-cli-clean-code

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/forge-cli-clean-code/proposal.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Remove dead code from run.go | pending | tasks/1-remove-dead-code.md |
| 2 | Add build artifacts to .gitignore | pending | tasks/2-gitignore-artifacts.md |
| 3 | Split forensic.go into functional files | pending | tasks/3-split-forensic.md |
| 4 | Split worktree.go into command files | pending | tasks/4-split-worktree.md |
| 5 | Extend ParseFrontmatter as shared YAML parser | pending | tasks/5-unify-frontmatter-parser.md |
| 6 | Unify dependency check logic | pending | tasks/6-unify-dep-check.md |
| 7 | Extract defaultRunClaude to shared location | pending | tasks/7-extract-default-run-claude.md |
| 8 | Complete SetFeature migration and remove deprecated function | pending | tasks/8-migrate-setfeature.md |
| 9 | Clean up re-export layer in errors.go and output.go | pending | tasks/9-clean-reexport-layer.md |
| 10 | Refactor askAutoBehavior to data-driven loop | pending | tasks/10-refactor-ask-auto-behavior.md |
| 11 | Fix validateRecordData os.Exit to return error | pending | tasks/11-fix-validate-record-data.md |
| 12 | Extract runE2ERegression to reduce nesting | pending | tasks/12-extract-e2e-regression.md |
| 13 | Unify error handling pattern across commands | pending | tasks/13-unify-error-handling.md |
| 14 | Migrate testbridge underlying functions to pkg/task | pending | tasks/14-testbridge-cleanup.md |
