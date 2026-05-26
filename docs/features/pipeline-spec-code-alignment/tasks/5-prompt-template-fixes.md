---
id: "5"
title: "Fix prompt template references (test-run.md and SCOPE labels)"
priority: "P0"
estimated_time: "1.5h"
dependencies: [4]
surface-key: "cli"
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: Fix prompt template references (test-run.md and SCOPE labels)

## Description

Fix two critical issues in prompt data templates embedded in the Go binary:

1. **test-run.md** (`forge-cli/pkg/prompt/data/test-run.md` lines 11, 31): References `forge:run-e2e-tests` skill which **does not exist**. The actual skill is `forge:run-tests`. This causes test.run tasks to fail immediately on execution. **Critical fix.**

2. **SCOPE label** (17 files in `forge-cli/pkg/prompt/data/*.md`): Templates use `SCOPE: {{SURFACE_KEY}}` tag. The label name `SCOPE` misleads developers into thinking it uses the deprecated scope field. Rename to `SURFACE_KEY:` or `SURFACE:` for clarity.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence G8 (test-run.md references non-existent forge:run-e2e-tests) and G10 (SCOPE misleading label in 17 templates)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Urgency` — H17 is Critical: test.run tasks fail directly
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for test-run.md fix and SCOPE label rename

## Acceptance Criteria
- [ ] `test-run.md` references `forge:run-tests` (not `forge:run-e2e-tests`)
- [ ] All 17 prompt templates use `SURFACE_KEY:` (or `SURFACE:`) label instead of `SCOPE:`
- [ ] Go code that injects/replaces these labels is updated to match
- [ ] Existing tests pass (`go test ./...`)

## Hard Rules
- Must update both the template .md files AND the Go code that processes these labels (prompt.go)
- Keep the `{{SURFACE_KEY}}` variable reference unchanged — only rename the label prefix

## Implementation Notes
- The SCOPE label appears in prompt.go where `TASK_CATEGORY` and other replacements happen. Check if there's a hardcoded `SCOPE:` string that also needs updating.
- For test-run.md, the fix is a simple string replacement: `forge:run-e2e-tests` → `forge:run-tests`.
