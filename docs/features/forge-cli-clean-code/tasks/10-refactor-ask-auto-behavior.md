---
id: "10"
title: "Refactor askAutoBehavior to data-driven loop"
priority: "P1"
estimated_time: "1h"
dependencies: [9]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 10: Refactor askAutoBehavior to data-driven loop

## Description
`askAutoBehavior()` is 130 lines with 13 identical `askConfirm` blocks. Refactor to a data-driven loop with a slice of prompt configs, reducing to <30 lines. Phase 4 (anti-pattern fix).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/config.go` — Contains `askAutoBehavior()` (search for exact location)

## Acceptance Criteria
- [ ] `askAutoBehavior()` reduced to <30 lines
- [ ] All 13 prompt behaviors preserved with identical semantics
- [ ] Data-driven approach: prompt configs defined as a slice, loop iterates over them
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- Behavioral equivalence: each prompt must ask the same question with the same defaults as before

## Implementation Notes
- The 13 identical blocks suggest a common pattern: extract a struct for prompt config (question text, default value, config key)
- Verify no block has unique logic that breaks the pattern
