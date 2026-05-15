---
id: "3"
title: "ui-design SKILL.md and template TUI support"
priority: "P1"
estimated_time: "2h"
dependencies: ["1", "2"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 3: ui-design SKILL.md and template TUI support

## Description
The core task: modify `ui-design/SKILL.md` to add TUI platform branch logic, and modify the `ui-design.md` output template to include the TUI component template with all 5 structural requirements from the lesson.

When platform=tui, the skill should: present TUI theme selection (Modern Dark / Minimal ASCII / DESIGN.md custom), use the TUI component template for each panel, and output `ui-design-tui.md` (or `ui-design-web.md` + `ui-design-tui.md` for multi-platform features).

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D1, D5, D7 sections)
- `docs/lessons/lesson-tui-tech-design-mockup.md` — 5 structural requirements to embed in template
- `plugins/forge/skills/ui-design/SKILL.md` — Current skill logic to modify
- `plugins/forge/skills/ui-design/templates/ui-design.md` — Current output template to modify
- `plugins/forge/skills/ui-design/templates/platforms/tui.md` — Created by Task 1
- `plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md` — Created by Task 1
- `plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md` — Created by Task 1

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/ui-design/SKILL.md` | Add TUI platform detection, TUI theme selection flow, TUI-specific output logic, multi-platform file splitting (proposal D7) |
| `plugins/forge/skills/ui-design/templates/ui-design.md` | Add TUI component template section: Panel Placement, ASCII Layout Mockup, Dimensions, Character Palette, Color Mapping, Edge Cases, States, Key Bindings, Data Binding (proposal D5) |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] `SKILL.md` detects platform=tui from PRD and enters TUI branch
- [ ] TUI branch presents theme selection: Modern Dark, Minimal ASCII, or DESIGN.md custom (proposal D3)
- [ ] TUI branch uses `platforms/tui.md` for navigation rules and selected theme for visual style
- [ ] `ui-design.md` template includes TUI component template with all 5 structural requirements from lesson: ASCII Layout Mockup, Dimensions (concrete numbers), Character Palette (Unicode + reason), Color Mapping (fg/bg color codes), Edge Cases (5 mandatory scenarios)
- [ ] Multi-platform features (e.g., web + tui) produce separate files: `ui-design-web.md` + `ui-design-tui.md` (proposal D7)
- [ ] Single TUI feature produces `ui-design-tui.md`
- [ ] Existing web/mobile behavior unchanged

## Hard Rules
- The 5 structural requirements in the TUI component template are MANDATORY (not optional). Each panel MUST include all 5 items. This is the key lesson from `deep-drill-analytics` — without enforcement, agents skip visual specs.

## Implementation Notes
- Study how SKILL.md currently handles web vs mobile platform branching — TUI should follow the same pattern
- The TUI component template (proposal D5) defines the exact section structure for each panel — use it as-is
- For multi-platform: single run reads PRD once, then splits output per platform. Each platform gets its own ui-design file and prototype directory (proposal D7)
- Character Palette must specify exact Unicode characters with code points, not vague descriptions — this is what prevents the "iterative trial-and-error" problem from the lesson
- Dimensions must be concrete numbers (e.g., "panel width: 60 chars"), not fuzzy descriptions (e.g., "takes up most of the screen")
