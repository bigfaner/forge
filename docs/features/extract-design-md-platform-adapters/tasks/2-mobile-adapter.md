---
id: "2"
title: "Implement mobile adapter and DESIGN.md template"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Implement mobile adapter and DESIGN.md template

## Description

Add mobile platform extraction to `extract-design-md`. Mobile adapter reuses the existing web CSS extraction pipeline (Layers 1-5) with a mobile User-Agent viewport context, then adds mobile-specific analysis: responsive breakpoints, touch target sizes, and safe area handling.

## Reference Files
- `docs/proposals/extract-design-md-platform-adapters/proposal.md` — Source proposal
- `plugins/forge/commands/extract-design-md.md` — Command to modify
- `plugins/forge/skills/ui-design/templates/styles/` — Reference for mobile-compatible style structures

## Acceptance Criteria
- [ ] `--platform mobile <url>` fetches the URL with mobile User-Agent context
- [ ] Responsive breakpoint analysis extracts common breakpoints (320/375/414/768px)
- [ ] Touch target estimation analyzes interactive element sizes from CSS (minimum 44x44pt guideline)
- [ ] Safe area handling inference from CSS `env(safe-area-inset-*)` or viewport meta tags
- [ ] Mobile DESIGN.md output extends the web template with additional sections: Touch Targets, Safe Areas, Responsive Breakpoints
- [ ] Match strategy still works: "match closest built-in style" or "fully custom" both produce mobile-flavored output
- [ ] Output is consumable by `/ui-design` without modification

## Hard Rules
- Reuse existing WebFetch infrastructure — no new tool dependencies
- Mobile extraction depends on the target URL serving responsive CSS — document this limitation clearly in the command

## Implementation Notes
- Mobile adapter is a pragmatic extension: same CSS extraction layers with mobile viewport headers
- If CSS lacks responsive tokens, fallback to Layer 5 (visual inference) for missing mobile tokens
- All mobile-specific values extracted from CSS should be clearly labeled; estimated values marked `(estimated)`
- The mobile DESIGN.md template should be an extension of the web template, not a completely separate structure
