---
id: "2"
title: "Bump patch version"
priority: "P2"
estimated_time: "5m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Bump patch version

## Description
Bump the forge CLI patch version to reflect the bug fix for task claim priority.

## Reference Files
- `docs/proposals/fix-task-claim-priority/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] Patch version bumped in plugin manifest and market config
- [ ] `forge --version` reports the new version

## Implementation Notes
- Use `/upgrade-forge` skill or manually bump version in manifest files
