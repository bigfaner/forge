---
id: "6"
title: "Bump CLI version and verify prompt generation includes principles"
priority: "P2"
estimated_time: "30m"
dependencies: [1, 2, 3, 4, 5]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: Bump CLI version and verify prompt generation includes principles

## Description

After all 5 template modifications are complete, bump the CLI patch version and verify that `forge prompt get-by-task-id` correctly includes the `<CODING_PRINCIPLES>` text for each coding type.

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal

## Acceptance Criteria
- CLI version patch-bumped (check `forge-cli/CLAUDE.md` section "Version Bump" for mechanism — currently references `scripts/version.txt` which may not exist; fall back to ldflags in justfile/makefile)
- Forge CLI rebuilt: `just compile backend`
- For each coding type (feature, enhancement, fix, refactor, cleanup), `forge prompt get-by-task-id <any-valid-id>` output contains `<CODING_PRINCIPLES>` text block
- Existing prompt-related tests pass: `go test -race -cover ./forge-cli/pkg/prompt/...`

## Hard Rules
- MUST NOT modify any template content — this task only bumps version and verifies
- MUST run `go test -race ./forge-cli/...` to catch regressions from template changes

## Implementation Notes
- The version mechanism may be ldflags-based (`forge-cli/pkg/version/version.go`: `var Version = "dev"`). Check the justfile for the actual bump workflow.
- If `scripts/version.txt` doesn't exist, check if the justfile/build script has a version bump target. If neither exists, update `version.go` directly and document the convention gap.
- Verification: pick any existing task ID (or create a temporary one) and run `forge prompt get-by-task-id` for each template type to confirm principles appear in the generated prompt.
