---
id: "5"
title: "Remove Version Bump rules from forge-cli/CLAUDE.md"
priority: "P1"
estimated_time: "0.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 5: Remove Version Bump rules from forge-cli/CLAUDE.md

## Description

Remove the Version Bump section (lines 28-33) from `forge-cli/CLAUDE.md`. Version management is now centralized in the `/release-cli` command (Task 4) — developers bump version when preparing a release, not on every code change. Daily code changes should no longer require modifying `version.txt`.

## Reference Files

- `forge-cli/CLAUDE.md`: contains the Version Bump section to remove at lines 28-33 (source: proposal.md#Implementation-5)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) |  |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/CLAUDE.md` | Remove the Version Bump section (lines 28-33) |

### Delete
| File | Reason |
|------|--------|
| (none) |  |

## Acceptance Criteria

- [ ] Version Bump section (lines 28-33, containing "Code changes must bump the version..." and the 3 semver bullet points) is removed from `forge-cli/CLAUDE.md`
- [ ] No other content in `forge-cli/CLAUDE.md` is modified

## Implementation Notes

- The removed content:
  ```
  ### Version Bump
  Code changes must bump the version in `scripts/version.txt`. Follow semver:
  - Patch: bug fixes, dead code removal (x.y.Z)
  - Minor: new features, new commands (x.Y.z)
  - Major: breaking CLI changes (X.y.z)
  ```
