---
id: "1"
title: "Rewrite run-tasks.md for token efficiency"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "documentation"
mainSession: false
noTest: true
scope: "backend"
---

# 1: Rewrite run-tasks.md for token efficiency

## Description

Compress `plugins/forge/commands/run-tasks.md` from ~250 lines to ~150 lines while preserving all semantics. Three changes: (1) replace `forge task query` with `forge task status` in Step 2b, (2) silent Breaking Gate and E2E Gate execution, (3) remove verbose explanations and condense error handling.

## Reference Files
- `docs/proposals/slim-run-tasks/proposal.md` — Source proposal
- `plugins/forge/commands/run-tasks.md` — Target file (current version)
- `forge-cli/internal/cmd/status.go` — `forge task status` query mode (outputs TASK_ID + STATUS only)
- `forge-cli/internal/cmd/query.go` — `forge task query` (outputs TASK_ID + STATUS + SCOPE + BREAKING, being replaced)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Compress from ~250 to ~150 lines; replace query→status; silent gates; condense error handling |

## Acceptance Criteria

- [ ] `run-tasks.md` reduced from ~250 to ~150 lines (preserve mermaid diagram, EXTREMELY-IMPORTANT blocks verbatim)
- [ ] Step 2b uses `forge task status <TASK_ID>` instead of `forge task query <TASK_ID>`
- [ ] Breaking Gate (Step 3a): test output redirected to file, only exit code checked on success; on failure, tail last 20 lines for diagnostic context before creating fix-task
- [ ] E2E Gate (Step 3b): same silent treatment — redirect output, check exit code, tail on failure
- [ ] Error Handling table condensed to inline notes or compact format
- [ ] All functionality preserved: claim, dispatch+verify, fix-task spawning, main-session routing, record-missing recovery, post-completion summary
- [ ] Pre-flight checks for Breaking Gate preserved (justfile existence, test recipe check)
- [ ] Pre-flight checks for E2E Gate preserved (test-e2e recipe, spec directory check)

## Hard Rules

- MUST keep mermaid flowchart diagram
- MUST preserve all `<EXTREMELY-IMPORTANT>` blocks verbatim
- MUST preserve Scope Resolution protocol reference
- MUST NOT change CLI commands (only changes the skill document)

## Implementation Notes

**Verification → status**: `forge task status <TASK_ID>` (query mode, 1 arg) outputs only TASK_ID and STATUS — exactly what Step 2b needs. The SCOPE and BREAKING fields from `query` are already captured in Step 1 claim output and don't need re-fetching.

**Silent execution pattern** for Breaking Gate:
```bash
just test [scope] > .forge/tmp/test-output.txt 2>&1; TEST_EXIT=$?
if [ $TEST_EXIT -ne 0 ]; then
  tail -20 .forge/tmp/test-output.txt
  # ... create fix-task with context from tail
fi
```

**Silent execution pattern** for E2E Gate:
```bash
just test-e2e --feature "$FEATURE" > .forge/tmp/e2e-output.txt 2>&1; E2E_EXIT=$?
if [ $E2E_EXIT -ne 0 ]; then
  tail -20 .forge/tmp/e2e-output.txt
  # ... create fix-task with context from tail
fi
```

**Compression targets**: verbose output parsing instructions in Step 1 (claim already documents its fields), Step 2b explanatory text, Step 3a/3b code examples (merge pre-flight + execution), Error Handling table (fold into Step sections), Post-Completion section (simplify).
