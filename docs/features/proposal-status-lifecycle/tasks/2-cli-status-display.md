---
id: "2"
title: "Update Go CLI display for Approved/Completed status"
priority: "P2"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 2: Update Go CLI display for Approved/Completed status

## Description

Verify that `forge proposal list` and `forge feature status` correctly display the new `Approved` and `Completed` status values. The current CLI reads the status string from frontmatter and displays it as-is, so the primary work is verifying existing behavior and adding any needed display improvements.

## Reference Files
- `docs/proposals/proposal-status-lifecycle/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `forge proposal list` displays proposals with `status: Approved` showing "Approved" in the STATUS column
- [ ] `forge proposal list` displays proposals with `status: Completed` showing "Completed" in the STATUS column
- [ ] `forge feature status <slug>` correctly reflects when a feature's manifest status is "completed"
- [ ] Existing tests pass (`go test ./...`)
- [ ] New test cases cover Approved and Completed status display

## Hard Rules
- Follow TDD: write failing tests first, then implement
- Run quality gate: `go build ./...` → `go vet ./...` → `golangci-lint run ./...` → `go test -race -cover ./...`
- Bump version in `scripts/version.txt` (patch bump for display enhancement)

## Implementation Notes
- Key files: `forge-cli/internal/cmd/proposal.go` (list + detail display), `forge-cli/internal/cmd/feature.go` (feature status display), `forge-cli/pkg/proposal/proposal.go` (proposal model)
- The Status field is a raw string — no validation or enum exists. The CLI displays whatever is in the frontmatter. Verify this works for "Approved" and "Completed", then add tests confirming it.
- Check if `forge feature status` needs any update — it reads manifest status (which uses "completed" lowercase) not proposal status.
