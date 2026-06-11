---
id: "4"
title: "TUI HTML prototype simulation rules"
priority: "P2"
estimated_time: "1h"
dependencies: ["3"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 4: TUI HTML prototype simulation rules

## Description
Modify the `prototype.md` template to add rules for generating HTML terminal simulation prototypes. TUI prototypes use HTML + CSS to simulate terminal appearance (black background, monospace font) for human review in the browser. The HTML is a visual approximation — ASCII mockups in ui-design.md remain the precise specification for agents.

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D2, D6 sections)
- `plugins/forge/skills/ui-design/templates/prototype.md` — Current prototype template
- `plugins/forge/skills/ui-design/templates/ui-design.md` — Modified by Task 3 (TUI component template)

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/ui-design/templates/prototype.md` | Add TUI prototype rules: single index.html with terminal-window div, all panels rendered, simulated key buttons for panel switching (proposal D6) |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] `prototype.md` includes TUI-specific prototype generation rules
- [ ] TUI prototype is a single `index.html` with all panels rendered inside a black "terminal window" div (proposal D6)
- [ ] Includes simulated key buttons at the bottom: `[Tab]`, `[1]`, `[q]`, `[:command]` to switch between panels
- [ ] Uses monospace font and dark background to approximate terminal appearance
- [ ] Panel layout in HTML matches the ASCII mockup from ui-design.md
- [ ] TUI prototypes output to `prototype/tui/` (multi-platform) or `prototype/` (single TUI feature)

## Implementation Notes
- Study how `prototype.md` currently specifies web/mobile prototype generation — TUI rules should follow the same pattern
- The HTML prototype is a human review tool, not an agent specification. The ASCII mockup + numeric dimensions in ui-design.md are the precise spec for implementation
- Proposal D2 explains the rationale: reuses existing HTML prototype infrastructure, browser viewing is intuitive, "distortion" is bounded by the precise numeric specs
- CSS should simulate common terminal characteristics: dark background (#1e1e1e or similar), monospace font (Menlo/Consolas/etc), fixed-width characters, no anti-aliasing effects
