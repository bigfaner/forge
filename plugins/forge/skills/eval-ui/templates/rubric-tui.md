# TUI Design Evaluation Rubric

**Total: 1000 points**
**Platform: tui (Terminal User Interface)**
**Report template:** `plugins/forge/skills/eval-ui/templates/report.md`

## Perspectives

Each dimension represents an independent stakeholder perspective. The scorer must evaluate from that role's standpoint — not from a generic "quality" viewpoint.

| Perspective | Role | Core Question |
|-------------|------|---------------|
| Requirement Coverage | Product Manager | Are all PRD UI requirements covered? Edge cases? |
| Terminal Experience | End User | Is it efficient for keyboard-driven terminal usage? |
| Visual Specification | Designer | Are ASCII mockups, character palettes, and color mappings precise? |
| Implementability | Developer | Can I code from this without guessing? |

## Required Sections

Mark missing required sections as 0 pts for the relevant dimension:

| Section | Required |
|---------|----------|
| ASCII Layout Mockup (per panel) | yes |
| Dimensions (concrete numeric values) | yes |
| Character Palette (Unicode + rationale per element) | yes |
| Color Mapping (foreground/background color codes) | yes |
| Edge Cases (5 mandatory scenarios per panel) | yes |
| Key Bindings table | yes |
| Data Binding table | yes |

## Dimensions

### 1. Requirement Coverage (250 pts) — Product Manager Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI function coverage | 0-80 | Does every UI function from `prd-ui-functions.md` have a corresponding panel in the design? Any gaps? |
| Navigation coverage | 0-40 | Are all keymaps, panel transitions, and mode switches covered? If PRD defines `## Navigation Architecture`, does the design cover all Keymap entries and Panel Layout definitions? |
| State coverage | 0-80 | Are all states (loading, empty, error, populated) addressed for each panel? Are mode-specific states defined (e.g., normal mode vs. command mode)? |
| Edge case handling | 0-50 | Are boundary conditions addressed: content wider than terminal width, long text overflow, no data, permission denied, terminal resize, concurrent actions? Or does the design only show the happy path? |

### 2. Terminal Experience (250 pts) — End User Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Keyboard reachability | 0-70 | Are all actions reachable via keyboard? Are keybindings ergonomic (home row preferred)? Is there a consistent keybinding scheme across panels? |
| Information scanning efficiency | 0-60 | Can a user scan the terminal output and immediately understand what matters most? Is information density appropriate — not too sparse, not overwhelming? |
| Panel switching intuitiveness | 0-60 | Is panel navigation intuitive (e.g., Tab/number keys for panel switch)? Can the user predict which panel they will land on? |
| Mode switching consistency | 0-60 | Are mode switches (e.g., normal → command → insert) consistent? Is the current mode always visible? Is there a clear visual indicator of the active mode? |

### 3. Visual Specification (250 pts) — Designer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| ASCII mockup completeness | 0-70 | Does every panel have a complete ASCII layout mockup using box-drawing characters? Is the layout structure clear with borders, dividers, and alignment? Missing ASCII mockup for any panel → this criterion = 0 for that panel. |
| Character palette precision | 0-60 | Is every visual element mapped to a specific Unicode character with rationale? Are box-drawing characters (├┤┬┴┼─│), block elements (▄▪), and other characters specified with their code points? Or are there ambiguous choices like "some line character"? |
| Color mapping compliance | 0-60 | Is every colored element mapped to specific terminal color codes (16-color or 256-color)? Are foreground and background colors both specified? Is the chosen color space (16-color/256-color) consistent with the declared theme? |
| Dimension specificity | 0-60 | Are all sizes stated as concrete numbers (e.g., "panel width: 60 chars", "status bar height: 1 line")? Or are there vague descriptions like "fits the screen" or "appropriate size"? |

### 4. Implementability (250 pts) — Developer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Dimension precision | 0-80 | Are all panel widths, heights, margins, and padding stated as exact character counts? Can a developer set layout parameters without guessing? |
| Mandatory edge case coverage | 0-80 | Does each panel cover the 5 mandatory edge case scenarios: (1) content exceeds panel width, (2) content exceeds panel height, (3) empty data, (4) terminal resize, (5) long text/overflow? Missing any → -50 per missing scenario. |
| Character Unicode clarity | 0-90 | Is every character specified with its Unicode code point (e.g., U+2500 for ─)? Or are characters described only by appearance ("horizontal line")? Can a developer copy-paste the exact character from the spec? |

## Deduction Rules

- **Missing ASCII mockup for a panel**: Visual Specification criterion for that panel = 0 pts
- **"待定" (pending) characters**: -30 pts per instance (a character must be specified, not deferred)
- **Missing mandatory edge case**: -50 pts per missing edge case scenario
- **Vague dimension descriptions** ("appropriate size", "fits screen"): -20 pts per instance
- **Missing required section**: 0 pts for that dimension
- **Cross-section inconsistency**: -30 pts per conflict (e.g., keybinding contradicts panel layout)
- **Happy-path only design** (no error/empty/loading states): -50 pts from Terminal Experience
- **Navigation Architecture gap**: -20 pts per PRD navigation entry not covered in design (from Requirement Coverage)
- **PRD UI function gap**: -30 pts per unaddressed UI function (from Requirement Coverage)
- **Orphan UI elements** (no data binding): -30 pts per element (from Implementability)
- **Placeholder text ("TBD", "TODO", "lorem ipsum")**: -20 pts per instance
