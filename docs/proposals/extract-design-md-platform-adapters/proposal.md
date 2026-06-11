---
created: 2026-05-16
author: "faner"
status: Draft
---

# Proposal: Extract-Design-MD Platform Adapters

## Problem

`/extract-design-md` only supports web applications (CSS extraction from URL), while its downstream consumer `/ui-design` supports three platforms: web, mobile, and TUI. This asymmetry forces users to manually create DESIGN.md for non-web platforms instead of extracting from existing apps.

### Evidence

- The skill description says "Analyze a web app's visual style" — explicitly web-only.
- TUI has two built-in themes (`modern-dark-tui`, `minimal-ascii`) but no way to extract from an existing terminal app.
- Mobile has no extraction path at all — users must adapt web results manually.
- `todo.txt` item 66: "optimize extract-design-md: adapt to different UI types" — acknowledged gap.

### Urgency

Low. This is a consistency improvement, not a blocking issue. Users can work around it by manually writing DESIGN.md or using built-in TUI themes. But the gap creates friction when onboarding TUI or mobile projects into the forge pipeline.

## Proposed Solution

Add platform adapters to `extract-design-md` with an explicit `--platform` flag (web/mobile/tui). Each adapter handles platform-specific input parsing, extraction logic, and output template. Web keeps existing behavior as default.

- **Web** (default): Unchanged 5-layer CSS extraction from URL.
- **Mobile**: Same URL-based extraction with mobile User-Agent + additional analysis of responsive breakpoints, touch target sizes, and safe area handling.
- **TUI**: Screenshot-based visual analysis using AI vision to infer ANSI color palette, character set, panel layout dimensions, and key bindings.

Each platform uses a dedicated DESIGN.md output template matching the corresponding ui-design style structure.

### Innovation Highlights

Straightforward adapter pattern. The creative insight is the TUI extraction approach: using AI vision to reverse-engineer terminal design tokens from a screenshot — there's no CSS to parse, so the entire extraction relies on visual inference. All TUI values are marked `(estimated)` to set expectations.

Mobile extraction is a pragmatic extension: reusing the web extraction pipeline with mobile viewport context rather than building a separate native app analysis path.

## Requirements Analysis

### Key Scenarios

- **Web (existing)**: User provides URL → CSS extraction → DESIGN.md. No behavior change.
- **Mobile web**: User provides URL with `--platform mobile` → mobile User-Agent fetch → responsive breakpoint analysis + touch target estimation → mobile DESIGN.md with safe area and touch target sections.
- **TUI screenshot**: User provides screenshot path with `--platform tui` → AI vision analysis → ANSI color palette + character set inference → TUI DESIGN.md matching modern-dark-tui structure.
- **Match strategy**: All platforms support "match closest built-in" or "fully custom". TUI matches against modern-dark-tui/minimal-ascii; mobile matches against web styles with mobile-specific overrides.

### Non-Functional Requirements

- TUI extraction quality depends on AI vision accuracy — all values must be marked `(estimated)`.
- Mobile extraction reuses existing WebFetch infrastructure — no new tool dependencies.
- Command remains a single file (command, not skill directory) — consistent with its utility classification.

### Constraints & Dependencies

- `allowed_tools` must be updated: TUI mode needs image analysis capability (Read for images or vision tools).
- TUI screenshot must be provided as a local file path (not URL).
- Mobile extraction depends on the target URL serving responsive CSS.

## Alternatives & Industry Benchmarking

### Industry Solutions

Design token extraction tools (Figma Token Studio, Style Dictionary) typically target web only. Terminal UI theming (oh-my-posh, starship) uses static config files, not extraction. No known tool extracts design tokens from terminal screenshots.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Inconsistent with ui-design; manual work for TUI/mobile | Rejected: defeats the purpose of extract-design-md |
| Platform adapters | This proposal | Full parity with ui-design; clean separation per platform | More code; TUI accuracy limited by AI vision | **Selected: best consistency/effort ratio** |
| TUI-only | — | Smaller scope | Mobile gap remains; half-measure | Rejected: incomplete parity |

## Feasibility Assessment

### Technical Feasibility

All extraction methods use existing Claude capabilities:
- Web/mobile: WebFetch (already in `allowed_tools`)
- TUI: Read tool for image analysis (already supports image files)
- AI vision for TUI screenshot analysis is built into Claude

### Resource & Timeline

Small scope — modifying one command file (~270 lines). Estimated 4-6 tasks.

### Dependency Readiness

No external dependencies. All tools already available in the Claude environment.

## Scope

### In Scope

- `--platform` argument (web/mobile/tui) with web as default
- Mobile adapter: mobile User-Agent fetch, responsive breakpoint extraction, touch target estimation, safe area inference
- TUI adapter: screenshot-based visual analysis (ANSI colors, character set, panel layout, key bindings)
- Platform-specific DESIGN.md output templates (mobile extends web template; TUI uses separate structure)
- Updated command description and argument-hints
- Updated `allowed_tools` if needed for TUI image analysis

### Out of Scope

- Multi-platform batch extraction (run once for all platforms)
- TUI source code parsing (only screenshots)
- Native mobile app analysis (only mobile web URLs)
- Changes to `/ui-design` skill (consumer remains unchanged)
- Changes to eval rubrics

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| TUI visual inference inaccuracy | High | Medium | Mark all TUI values as `(estimated)`; encourage manual review |
| Mobile CSS lacks responsive tokens | Medium | Low | Fallback to visual inference layer (Layer 5) for missing mobile tokens |
| Command file grows too large | Medium | Low | Keep adapter logic concise; shared match-strategy step remains unified |
| TUI screenshot quality varies | Medium | Medium | Require clear screenshot guidelines; reject blurry/low-res inputs early |

## Success Criteria

- [ ] `/extract-design-md --platform tui <screenshot>` generates a TUI DESIGN.md with ANSI color palette, character palette, panel layout, and key bindings
- [ ] `/extract-design-md --platform mobile <url>` generates a mobile DESIGN.md with touch targets, safe areas, and responsive breakpoints
- [ ] `/extract-design-md <url>` (no flag) produces identical output to current behavior
- [ ] All platform outputs are consumable by `/ui-design` without modification
- [ ] TUI output matches the structure of built-in TUI themes (modern-dark-tui/minimal-ascii)
