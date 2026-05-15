---
id: "1"
title: "TUI platform definition and themes"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 1: TUI platform definition and themes

## Description
Create the foundational TUI platform definition file and two TUI design themes. This is the base layer that ui-design SKILL.md and templates will reference when platform=tui.

The proposal adds TUI as a third platform alongside web and mobile in `/ui-design`. The platform definition specifies TUI navigation rules (Keymap + Panel Layout + Mode), and the themes define character sets, color spaces, and density for TUI UI design.

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D1, D3, D5 sections)
- `docs/lessons/lesson-tui-tech-design-mockup.md` — 5 structural requirements from past TUI experience
- `plugins/forge/skills/ui-design/templates/platforms/web.md` — Reference: existing web platform definition
- `plugins/forge/skills/ui-design/templates/platforms/mobile.md` — Reference: existing mobile platform definition
- `plugins/forge/skills/ui-design/templates/styles/apple.md` — Reference: existing style file format

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/ui-design/templates/platforms/tui.md` | TUI platform navigation rules (Keymap, Panel Layout, Mode, Navigation Rules) |
| `plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md` | Modern Dark theme: 256-color, box-drawing chars, dark background, compact density |
| `plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md` | Minimal ASCII theme: 16-color, pure ASCII chars, default terminal bg, loose density |

### Modify
| File | Changes |
|------|---------|
| — | — |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] `platforms/tui.md` defines TUI navigation structure: Keymap table `[Key | Action | Context/Mode]`, Panel Layout table `[Panel | View | Position | Size Hint]`, Modes table `[Mode | Description | Default Keybindings]`, Navigation Rules
- [ ] `styles/modern-dark-tui.md` specifies: color space (256-color), character set (box-drawing + block elements with examples: ▄▪─│┃), palette (dark bg, high contrast, green/red/blue semantic colors), density (compact), applicable scenarios
- [ ] `styles/minimal-ascii-tui.md` specifies: color space (16-color), character set (pure ASCII: `#=-\|*+.`), palette (default terminal bg, distinguish by weight/spacing), density (loose), applicable scenarios
- [ ] Both theme files follow the same format as existing style files (e.g., `apple.md`, `shadcn.md`)

## Implementation Notes
- Study existing `platforms/web.md` and `platforms/mobile.md` to match the established format and structure
- The TUI platform definition must include the 5 structural requirements from lesson (`docs/lessons/lesson-tui-tech-design-mockup.md`) as mandatory sections: ASCII Layout Mockup, Dimensions, Character Palette, Color Mapping, Edge Cases
- Proposal D3 specifies exact character sets and color spaces for each theme — use those directly
- Proposal D4 defines the TUI Navigation Architecture structure — use as the basis for `platforms/tui.md`
