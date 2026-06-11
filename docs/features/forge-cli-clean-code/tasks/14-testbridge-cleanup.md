---
id: "14"
title: "Migrate testbridge underlying functions to pkg/task"
priority: "P1"
estimated_time: "1.5h"
dependencies: [9]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 14: Migrate testbridge underlying functions to pkg/task

## Description
`testbridge.go` exports 37 internal symbols for test access, but the file is in the main build. Migrate the underlying function implementations (e.g. `getTaskPhase()`) to `pkg/task/`, keeping testbridge as thin aliases. This ensures tests don't need call-site changes. Phase 4 (anti-pattern fix).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/testbridge.go` — 37 exported symbols
- `forge-cli/pkg/task/` — Target location for implementations

## Acceptance Criteria
- [ ] Underlying function implementations moved to `pkg/task/`
- [ ] `testbridge.go` retains only thin aliases (one-line wrappers or type aliases)
- [ ] All test call sites continue working without modification
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- Do NOT modify test call sites — the testbridge alias layer ensures backward compatibility
- Each migrated function must have a clear home in `pkg/task/` (not just dumped into a file)

## Implementation Notes
- The XY analysis in the proposal refined this: goal is test accessibility, not removing testbridge
- Identify which of the 37 symbols have real logic vs. already-thin wrappers
- Migrate only symbols with substantial logic; leave trivial aliases in testbridge
