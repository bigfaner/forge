---
id: "7"
title: "gen-journeys skill adaptation for surfaces"
priority: "P1"
estimated_time: "1h"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 7: gen-journeys skill adaptation for surfaces

## Description

Update gen-journeys skill to query `forge surfaces <path>` instead of independently detecting surface type. Rename rule files to use new naming convention. Update SKILL.md instructions.

## Reference Files
- `proposal.md#CLI-命令与退出码契约` — gen-journeys calling contract (exit code 0 = parse stdout, exit code 1 = prompt user)
- `proposal.md#统一命名规范` — webui→web, mobileui→mobile, tui→tui
- `proposal.md#Key-Risks` — new/old field coexistence risk, rule file rename sync

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-journeys/rules/surface-web.md` | Renamed from surface-webui.md |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/SKILL.md` | Replace surface detection section with `forge surfaces <path>` query; update `surface` field reference to `surfaces` |
| `plugins/forge/skills/gen-jneys/rules/surface-mobile.md` | Rename from surface-mobileui.md (if applicable) |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-journeys/rules/surface-webui.md` | Renamed to surface-web.md |

## Acceptance Criteria

- [ ] SKILL.md "Surface Detection" section rewritten: detect via `forge surfaces <path>` CLI call instead of scanning project
- [ ] SKILL.md uses exit code contract: exit 0 → parse stdout for surface type, exit 1 → prompt user to configure
- [ ] `surface-webui.md` renamed to `surface-web.md`
- [ ] `surface-mobileui.md` renamed to `surface-mobile.md` (if it exists with old naming)
- [ ] Internal naming references updated: `webui` → `web`, `mobileui` → `mobile`
- [ ] Rule file loading logic updated to match new filenames

## Hard Rules

- Do NOT change the rule file content structure — only rename and update naming references
- The skill must NOT independently scan project files for surface detection anymore

## Implementation Notes

- SKILL.md lines 19-80 (Surface Detection section) need major rewrite
- Current rule files: `surface-api.md`, `surface-cli.md`, `surface-tui.md` already use short names — no rename needed
- Only `surface-webui.md` → `surface-web.md` and potentially `surface-mobileui.md` → check actual filename
