---
id: "4"
title: "Create types/tui.md instruction file"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 4: Create types/tui.md instruction file

## Description

Extract TUI-specific generation logic from the monolithic `gen-test-scripts` SKILL.md into a dedicated `types/tui.md` type instruction file. TUI tests verify terminal full-screen applications (like vim/htop) — character-level output, keyboard navigation, screen state transitions.

Modeled after `gen-test-cases/types/tui.md` structure.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/types/tui.md` — Reference structure (gen-test-cases TUI type file)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Source of TUI-relevant content (general reconnaissance and Fact Table, adapted for TUI)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/tui.md` | TUI type instruction file with conventions frontmatter |

## Acceptance Criteria

- [ ] `plugins/forge/skills/gen-test-scripts/types/tui.md` exists
- [ ] Frontmatter declares `type: tui` and `conventions: [testing-tui.md]`
- [ ] Contains a **Reconnaissance Strategy** section with TUI-specific search patterns (grep terminal libraries — bubbletea, tcell, tview, termbox; screen rendering functions; key binding registrations)
- [ ] Contains a **Fact Table Required Keys** section listing minimum keys for TUI type (TUI binary name or entry point, key binding definitions)
- [ ] Contains a **Verification Method** section describing how to confirm the project exposes a TUI (grep terminal framework imports, detect full-screen rendering patterns)
- [ ] Contains a **Generation Patterns** section describing how TUI test cases translate to executable scripts (non-interactive stdin pipe for key sequences, terminal output capture, screen state assertions, exit code + output comparison)
- [ ] Contains a **TUI Antipattern Guards** section (interactive prompts that require real terminal, testing via source code inspection instead of runtime, sleep-based waits for screen transitions)
- [ ] At least 3 section headings are unique to this file
- [ ] TUI vs CLI disambiguation: TUI clears terminal and redraws (full-screen); CLI produces line-oriented output. Document this distinction in the reconnaissance strategy.

## Hard Rules

- TUI test scripts must use non-interactive execution (stdin pipe, not real terminal) — no interactive test modes
- TUI tests that require real terminal interaction should be explicitly marked as "manual only" with a skip rationale

## Implementation Notes

- The current SKILL.md has no dedicated TUI branch — TUI shares generic reconnaissance with CLI. This task must create TUI-specific content from scratch, guided by the TUI profile's capabilities.
- Reference the `go-test` profile which has `tui` capability — understand what go-test TUI templates expect
- Key insight from gen-test-cases TUI disambiguation: TUI clears terminal and redraws (full-screen like vim/htop). Interactive prompts (line-by-line Q&A) are CLI, not TUI.
