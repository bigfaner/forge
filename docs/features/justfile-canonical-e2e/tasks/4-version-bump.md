---
id: "4"
title: "Version bump to 3.10.0"
priority: "P2"
estimated_time: "5m"
dependencies: ["3"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Version bump to 3.10.0

## Description

Bump the forge CLI version from 3.9.0 to 3.10.0 to reflect the e2e delegation refactor. This is a breaking change for any code that relied on the removed manifest.yaml command fields (none found in codebase).

## Reference Files
- `docs/proposals/justfile-canonical-e2e/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/scripts/version.txt` | Change `3.9.0` to `3.10.0` |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `scripts/version.txt` contains `3.10.0`
- [ ] `forge --version` (after build) reports `3.10.0`

## Implementation Notes

- Single line change. No other files reference the version number directly — it is read by the build system from `scripts/version.txt`.
