---
id: "3"
title: "Implement TUI adapter and DESIGN.md template"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Implement TUI adapter and DESIGN.md template

## Description

Add TUI platform extraction to `extract-design-md`. Unlike web/mobile, TUI has no CSS to parse — the adapter uses AI vision to analyze a terminal screenshot and reverse-engineer design tokens: ANSI color palette, character set, panel layout dimensions, and key bindings.

## Reference Files
- `docs/proposals/extract-design-md-platform-adapters/proposal.md` — Source proposal
- `plugins/forge/commands/extract-design-md.md` — Command to modify
- `plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md` — TUI style reference structure
- `plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md` — TUI style reference structure

## Acceptance Criteria
- [ ] `--platform tui <screenshot-path>` accepts a local file path to a terminal screenshot
- [ ] AI vision analysis extracts: ANSI color palette (xterm-256 color numbers), character set (box-drawing, block elements), panel layout dimensions (rows, columns), and key bindings
- [ ] TUI DESIGN.md output matches the structure of built-in TUI themes (modern-dark-tui/minimal-ascii): Color Space, Character Set, Character Palette Reference, Color Palette, Panel Layout sections
- [ ] All TUI extracted values are marked `(estimated)` to set accuracy expectations
- [ ] Match strategy supports: "match closest built-in TUI theme" (modern-dark-tui or minimal-ascii) or "fully custom"
- [ ] Rejects blurry or low-resolution screenshots with a clear error message
- [ ] Output is consumable by `/ui-design` without modification

## Hard Rules
- TUI input must be a local file path (not URL) — screenshots cannot be fetched remotely
- All values MUST be marked `(estimated)` — AI vision inference is inherently approximate
- TUI output structure must align with `modern-dark-tui.md` sections for `/ui-design` compatibility

## Implementation Notes
- Key risk: TUI visual inference inaccuracy (High likelihood, Medium impact). Mitigation: all values marked `(estimated)` and users encouraged to manually review
- Key risk: TUI screenshot quality varies. Mitigation: add screenshot quality guidelines to the command and reject clearly unreadable inputs early
- The creative insight here is using AI vision to reverse-engineer terminal design tokens from a screenshot — there's no CSS or structured source to parse
- Reference the `modern-dark-tui.md` structure for the output template: Color Space, Character Set with Character Palette Reference table, Color Palette with xterm-256 numbers, Panel Layout
