---
id: "3"
title: "Create per-type templates (templates/*-test-cases.md)"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
type: "documentation"
mainSession: false
noTest: true
---

# 3: Create per-type templates (templates/*-test-cases.md)

## Description

Create 5 per-type template files under `plugins/forge/skills/gen-test-cases/templates/`. Each template contains only that type's section — frontmatter, TC placeholders, and type-specific traceability. This replaces the current single `templates/test-cases.md` which has sections for all 5 types.

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/templates/test-cases.md` — Current monolithic template

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-cases/templates/ui-test-cases.md` | UI-only template with frontmatter, TC section, traceability table |
| `plugins/forge/skills/gen-test-cases/templates/tui-test-cases.md` | TUI-only template |
| `plugins/forge/skills/gen-test-cases/templates/mobile-test-cases.md` | Mobile-only template |
| `plugins/forge/skills/gen-test-cases/templates/api-test-cases.md` | API-only template |
| `plugins/forge/skills/gen-test-cases/templates/cli-test-cases.md` | CLI-only template |

### Modify
| File | Changes |
|------|---------|
| (none — keep existing `templates/test-cases.md` for backward compatibility) |

### Delete
| File | Reason |
|------|--------|
| (none — keep legacy template for backward compat) |

## Acceptance Criteria
- [ ] 5 template files created, one per type
- [ ] Each template has frontmatter with `feature`, `sources`, `generated` fields
- [ ] Each template has a single TC section for its type (no empty sections for other types)
- [ ] Each template has a traceability table scoped to its type
- [ ] Route Validation section present only in UI, TUI, Mobile templates
- [ ] Legacy `templates/test-cases.md` preserved unchanged for backward compatibility

## Hard Rules
- Templates must not include TCs from other types — single-type only
- The TC placeholder format matches the current format (TC-{NNN}, Source, Type, Target, Test ID, Pre-conditions, Route, Steps, Expected, Priority)

## Implementation Notes
- Extract each type's section from the current `templates/test-cases.md` into its own file
- The summary table in each template shows only that type's count (no multi-type rows)
- Route Validation section is only relevant for types that use Route fields (UI, TUI, Mobile)
