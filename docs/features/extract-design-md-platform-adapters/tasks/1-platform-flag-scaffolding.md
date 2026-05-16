---
id: "1"
title: "Add --platform flag and command scaffolding"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add --platform flag and command scaffolding

## Description

Update the `extract-design-md` command to accept a `--platform` argument (web/mobile/tui) with web as default. This establishes the routing layer that Tasks 2 and 3 will plug into.

## Reference Files
- `docs/proposals/extract-design-md-platform-adapters/proposal.md` — Source proposal
- `plugins/forge/commands/extract-design-md.md` — Command to modify (270 lines)

## Acceptance Criteria
- [ ] Command frontmatter updated: description mentions all three platforms, `allowed_tools` includes image analysis capability, `argument-hints` includes `--platform` with valid values (web/mobile/tui)
- [ ] Default behavior (`--platform web` or no flag) produces identical output to current behavior
- [ ] Platform routing logic added: when `--platform mobile` or `--platform tui` is provided, the command delegates to the appropriate adapter section (placeholder for Tasks 2/3)
- [ ] Input validation rejects invalid platform values with a clear error message

## Hard Rules
- Web extraction must remain byte-for-byte identical when `--platform web` (or no flag) is used — no behavioral drift
- Command must remain a single file (not converted to skill directory)

## Implementation Notes
- Add `--platform` to `argument-hints` with description: "Target platform: web (default), mobile, or tui"
- Update `allowed_tools` to include any image-related tools needed for TUI mode later
- The routing can be a simple conditional in the Process Flow section — no abstraction needed yet
- Key risk from proposal: command file growing too large — keep routing logic minimal
