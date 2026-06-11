---
id: "2"
title: "Create per-type instruction files (types/*.md)"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Create per-type instruction files (types/*.md)

## Description

Create 5 per-type instruction files under `plugins/forge/skills/gen-test-cases/types/`. Each file contains the type-specific Steps 3-4 instructions extracted from the current monolithic SKILL.md, including: type classification rules, priority assignment, TC format, target derivation rules, type-specific quality rules, and any type-specific generation logic (e.g., Integration TC generation for UI/Mobile, route validation specifics).

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Current monolithic (source of extracted content)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-cases/types/ui.md` | UI-specific Steps 3-4: page rendering, navigation, visual state, Route field, Integration TC generation for existing-page placements |
| `plugins/forge/skills/gen-test-cases/types/tui.md` | TUI-specific Steps 3-4: terminal screen rendering, keyboard navigation, output assertions, screen transitions |
| `plugins/forge/skills/gen-test-cases/types/mobile.md` | Mobile-specific Steps 3-4: touch interactions, gestures, screen transitions, accessibility labels, platform-specific components |
| `plugins/forge/skills/gen-test-cases/types/api.md` | API-specific Steps 3-4: endpoints, request/response, status codes, data contracts, HTTP methods |
| `plugins/forge/skills/gen-test-cases/types/cli.md` | CLI-specific Steps 3-4: commands, flags, output format, exit codes, arguments, stdin/stdout |

### Modify
| File | Changes |
|------|---------|
| (none) |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria
- [ ] 5 instruction files created: `types/ui.md`, `types/tui.md`, `types/mobile.md`, `types/api.md`, `types/cli.md`
- [ ] Each file has YAML frontmatter with `conventions` field listing type-specific convention dependencies
- [ ] Each file covers type-specific Steps 3-4 completely (classification rules, TC format, target derivation, quality rules)
- [ ] No content is lost from the current monolithic SKILL.md — all type-specific logic is preserved
- [ ] UI instruction file includes Integration TC generation for `existing-page` placements
- [ ] Convention frontmatter examples: UI → `conventions: [testing-ui.md, frontend.md]`, CLI → `conventions: [testing-cli.md]`

## Hard Rules
- Each instruction file must be self-contained for Steps 3-4 — an agent loading only the dispatcher + one type file must have everything needed to generate test cases for that type
- Do NOT duplicate type-agnostic content (Steps 0-2.5) that stays in the dispatcher

## Implementation Notes
- Extract content from SKILL.md lines 134-235 (Step 3: Classify & Generate, Integration TCs, Step 3.5 route validation specifics, Step 4 output specifics)
- Each type's classification indicators (the "Indicators" table rows) map to its instruction file
- Target derivation rules are per-type: UI→`ui/<page-name>`, TUI→`tui/<screen-name>`, Mobile→`mobile/<screen-name>`, API→`api/<resource>`, CLI→`cli/<command>`
- The Antipattern Prevention rules (6 rules) should be referenced in each instruction file but not duplicated — link back to the dispatcher's shared rules
