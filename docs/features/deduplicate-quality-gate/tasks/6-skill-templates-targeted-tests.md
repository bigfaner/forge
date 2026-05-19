---
id: "6"
title: "Update skill templates with targeted test instructions"
priority: "P2"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 6: Update Skill Templates with Targeted Tests

## Description

All coding task type templates (5 coding-*.md files) and fix/cleanup task templates (2 files) currently instruct agents to run `just test [scope]` during task execution. This is the project-wide test runner and is redundant with the CLI submit gate.

Replace `just test` instructions with **targeted test instructions**: framework-native commands on changed packages/files only. The agent runs fast, focused tests during development. The CLI submit gate handles full verification at submission time.

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (items 5, 6, Tier 1 Targeted Tests)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-feature.md` | Replace Step 3 `just test` with targeted test instructions |
| `forge-cli/pkg/prompt/data/coding-enhancement.md` | Replace Step 3 `just test` with targeted test instructions |
| `forge-cli/pkg/prompt/data/coding-cleanup.md` | Replace Step 3 `just test` with targeted test instructions |
| `forge-cli/pkg/prompt/data/coding-refactor.md` | Replace Step 3 `just test` with targeted test instructions |
| `forge-cli/pkg/prompt/data/coding-fix.md` | Replace Step 3 `just test` with targeted test instructions |
| `forge-cli/pkg/template/data/fix-task.md` | Replace `just test [scope]` with targeted test instructions |
| `forge-cli/pkg/template/data/cleanup-task.md` | Replace `just test [scope]` with targeted test instructions |

## Acceptance Criteria

- [ ] No template references `just test` for task-level verification (grep confirms zero matches in all 7 files)
- [ ] All 5 coding-*.md templates instruct agents to run targeted tests (e.g., `go test ./pkg/foo/...` for Go)
- [ ] fix-task.md and cleanup-task.md templates updated to targeted tests
- [ ] Templates still instruct agents to run `just compile`, `just fmt`, `just lint` for static checks
- [ ] Templates note that the CLI submit gate handles full verification (agents don't need `just test`)

## Implementation Notes

- **Coding templates** (5 files): Replace Step 3 (Full Verification / Quality Gate) with a targeted test step. Keep `just compile`, `just fmt`, `just lint` for static verification. Replace `just test {{SCOPE}}` with targeted test instructions like:
  - Go: `go test -race -cover ./changed/package/...`
  - Generic: run framework-native test command on affected files
  - Note: "Full project-wide tests run at CLI submit (forge task submit) — agent runs targeted tests only."
- **fix-task.md** and **cleanup-task.md**: Replace `just test [scope]` in Verification section with targeted test instructions. Keep the note about e2e regression being verified by the dispatcher.
- The `{{SCOPE}}` variable remains relevant for `just compile/fmt/lint` — only `just test` is replaced.
- These templates are embedded in the Go binary via `//go:embed`. Changes require rebuilding the CLI.
