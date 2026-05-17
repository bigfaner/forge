---
id: "1"
title: "Submit write-once protection"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: Submit write-once protection

## Description
Add record existence check in `forge task submit` before `os.WriteFile`. When a record file already exists, block the write unless `--force` is set. This prevents accidental record overwrites that destroy audit history.

The existing `--force` flag (currently used to bypass quality gate) gains a second semantic: allowing record overwrite. Both meanings align with "skip safety checks".

## Reference Files
- `docs/proposals/task-record-immutability/proposal.md` — Source proposal
- `forge-cli/internal/cmd/submit.go` — Primary implementation target (lines 164-170)
- `forge-cli/internal/cmd/submit_test.go` — Existing tests to extend

## Acceptance Criteria
- [ ] `forge task submit <id>` fails with descriptive error when record file already exists: "Record for task X already exists at <path>. Use --force to overwrite, or create a fix task instead."
- [ ] `forge task submit <id> --force` overwrites existing record with stderr warning: "WARNING: Overwriting existing record at <path>"
- [ ] `forge task submit <id>` succeeds normally when record file does not exist (current behavior preserved)
- [ ] Existing submit tests continue to pass

## Hard Rules
- Reuse existing `submitForce` variable — do NOT add a new flag
- Error uses `ErrValidation` category via `NewAIError`
- Warning goes to stderr, not stdout

## Implementation Notes
- Insert `os.Stat(recordPath)` check between `os.MkdirAll` (line 165) and `os.WriteFile` (line 168)
- If `os.Stat` returns nil (file exists): check `submitForce`, either error or warn+proceed
- If `os.IsNotExist(err)`: proceed to write (normal path)
- If other error from `os.Stat`: propagate as validation error
