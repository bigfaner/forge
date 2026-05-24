---
id: "2"
title: "Add build artifacts to .gitignore"
priority: "P0"
estimated_time: "15m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 2: Add build artifacts to .gitignore

## Description
Add residual build artifacts (cmd.out, cout.out, coverage.out, just.out) to the repository root `.gitignore`. These files should never have been tracked. This is Phase 1 (dead code elimination) of the cleanup proposal.

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `.gitignore` — Repository gitignore

## Acceptance Criteria
- [ ] Build artifact patterns added to root `.gitignore`
- [ ] Artifacts removed from git tracking (if currently tracked): `git rm --cached <files>`

## Hard Rules
- Only modify `.gitignore`, do not delete files from disk

## Implementation Notes
- Check which artifacts exist at repo root before deciding what to add
- Use `git rm --cached` to untrack without deleting local copies
