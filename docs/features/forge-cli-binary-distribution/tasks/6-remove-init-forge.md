---
id: "6"
title: "Remove /init-forge command"
priority: "P2"
estimated_time: "0.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 6: Remove /init-forge command

## Description

Delete `plugins/forge/commands/init-forge.md`. The installation flow is now handled by `install.sh` + `forge upgrade`, making the old compile-from-source command obsolete. The local developer build script (`forge-cli/scripts/install-local.sh`) is preserved.

## Reference Files

- `plugins/forge/commands/init-forge.md`: file to delete (source: proposal.md#Implementation-6)
- `docs/conventions/forge-distribution.md`: must load before modifying plugin files per CLAUDE.md mandatory rule (source: CLAUDE.md MANDATORY section)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) |  |

### Modify
| File | Changes |
|------|---------|
| (none) |  |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/commands/init-forge.md` | Obsoleted by install.sh + forge upgrade flow |

## Acceptance Criteria

- [ ] `plugins/forge/commands/init-forge.md` deleted
- [ ] No references to `/init-forge` remain in other plugin command files under `plugins/forge/commands/`

## Hard Rules

- Must load `docs/conventions/forge-distribution.md` before modifying any plugin files

## Implementation Notes

- `forge-cli/scripts/install-local.sh` is preserved — it serves a different purpose (developer local builds from source)
- Search for `/init-forge` references across `plugins/forge/` to ensure no dangling references
