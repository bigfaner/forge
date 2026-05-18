---
id: "6"
title: "Update forge-distribution.md and remove references/ directory"
priority: "P1"
estimated_time: "20m"
dependencies: ["1", "2", "3", "4", "5"]
type: "cleanup"
mainSession: false
---

# 6: Update forge-distribution.md and remove references/ directory

## Description
Update `docs/conventions/forge-distribution.md` to remove all documentation of the `references/` directory, then delete the entire `plugins/forge/references/` directory. This is the final cleanup task that runs after all content has been inlined and CLI files relocated.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — Distribution convention doc to update

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/forge-distribution.md` | Remove references/ from directory tree (lines 31-35), component table (line 54), path example (line 135), and reference file explanation (line 139) |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/references/` | Entire directory — all consumers updated in tasks 1-5 |

## Acceptance Criteria
- [ ] `plugins/forge/references/` directory no longer exists
- [ ] Zero occurrences of `references/shared/` in any file under `plugins/forge/`
- [ ] `docs/conventions/forge-distribution.md` has no mention of `references/` directory
- [ ] `docs/conventions/forge-distribution.md` directory tree and component table are accurate without references/
- [ ] Path resolution section no longer lists `references/shared/` cross-skill reference pattern

## Hard Rules
- Verify with `grep -r "references/shared/" plugins/forge/` that zero references remain before deleting the directory
- Do NOT remove the `references/` documentation from user project directory structure (line 86: `docs/reference/`) — that is a different concept

## Implementation Notes
- The distribution doc has 4 locations referencing `references/`:
  1. Directory tree listing (lines 31-35): remove the `references/` subtree
  2. Component table (line 54): remove the `references/` row
  3. Path rules table (line 135): remove the cross-skill example row
  4. Reference file explanation (line 139): remove the entire paragraph about `references/shared/`
- `profile-detection.md` has no consumers in the plugin — it can simply be deleted with the directory
