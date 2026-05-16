---
feature: "extract-design-md-platform-adapters"
sources:
  - docs/proposals/extract-design-md-platform-adapters/proposal.md
  - docs/features/extract-design-md-platform-adapters/tasks/1-platform-flag-scaffolding.md
  - docs/features/extract-design-md-platform-adapters/tasks/2-mobile-adapter.md
  - docs/features/extract-design-md-platform-adapters/tasks/3-tui-adapter.md
generated: "2026-05-16"
profile: "go-test"
---

# Test Cases: extract-design-md-platform-adapters

## Summary

| Type | Count |
|------|-------|
| CLI  | 18    |
| **Total** | **18** |

> **Profile**: go-test (capabilities: tui, api, cli). This feature is a CLI command, so all test cases are CLI type. No web UI or API interfaces exist.

---

## CLI Test Cases

### Platform Flag & Scaffolding

## TC-001: Default platform produces web output
- **Source**: Task 1 / AC-2, Proposal Success Criteria item 3
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/default-platform-produces-web-output
- **Pre-conditions**: A target URL is available (e.g. a local test server or known public URL). DESIGN.md does not already exist in the project root.
- **Steps**:
  1. Run `extract-design-md <url>` without any `--platform` flag
  2. Confirm the command completes successfully and writes DESIGN.md
  3. Verify DESIGN.md contains web-style sections: Color Palette (hex values), Typography, Components, Layout, Depth & Elevation
- **Expected**: Output is identical to the current web-only behavior. DESIGN.md contains standard web design tokens with no mobile-specific or TUI-specific sections.
- **Priority**: P0

## TC-002: Explicit web platform produces identical output to default
- **Source**: Task 1 / AC-2
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/explicit-web-platform-identical-output
- **Pre-conditions**: A target URL is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url>` (no flag) and capture output
  2. Delete DESIGN.md
  3. Run `extract-design-md <url> --platform web` and capture output
  4. Diff the two DESIGN.md files
- **Expected**: The two outputs are byte-for-byte identical. No behavioral drift when explicitly specifying `--platform web`.
- **Priority**: P0

## TC-003: Invalid platform value rejected with clear error
- **Source**: Task 1 / AC-4
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/invalid-platform-value-rejected
- **Pre-conditions**: None.
- **Steps**:
  1. Run `extract-design-md <url> --platform invalid`
  2. Observe the error output
  3. Verify DESIGN.md was not created
- **Expected**: Command outputs `ERROR: unsupported platform "invalid". Must be one of: web, mobile, tui` and exits without writing any file.
- **Priority**: P0

## TC-004: Command frontmatter describes all three platforms
- **Source**: Task 1 / AC-1
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/frontmatter-describes-all-platforms
- **Pre-conditions**: Command file exists at `plugins/forge/commands/extract-design-md.md`.
- **Steps**:
  1. Read the command file frontmatter
  2. Verify `description` mentions web, mobile, and TUI platforms
  3. Verify `argument-hints` includes `--platform` with valid values (web/mobile/tui)
  4. Verify `allowed_tools` includes Read (for image analysis in TUI mode)
- **Expected**: Command frontmatter is correctly updated with all platform references, the `--platform` argument hint, and image analysis capability in allowed_tools.
- **Priority**: P1

### Mobile Adapter

## TC-005: Mobile platform fetches with mobile context
- **Source**: Task 2 / AC-1
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/mobile-fetches-with-mobile-context
- **Pre-conditions**: A target URL serving responsive CSS is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. Verify the command completes successfully
  3. Verify DESIGN.md is written
- **Expected**: Command fetches the URL with mobile User-Agent/viewport context and generates DESIGN.md without errors.
- **Priority**: P0

## TC-006: Responsive breakpoint analysis extracts common breakpoints
- **Source**: Task 2 / AC-2
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/responsive-breakpoint-analysis
- **Pre-conditions**: A target URL with responsive CSS (containing `@media` queries) is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. Read the generated DESIGN.md
  3. Check for "Responsive Breakpoints" section
  4. Verify breakpoints for 320px, 375px, 414px, 768px are listed (or subset matching the site's CSS)
- **Expected**: DESIGN.md contains a "Responsive Breakpoints" section listing breakpoints extracted from the target site's CSS with layout change descriptions for each.
- **Priority**: P0

## TC-007: Touch target estimation analyzes interactive elements
- **Source**: Task 2 / AC-3
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/touch-target-estimation
- **Pre-conditions**: A target URL with interactive elements (buttons, links) is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. Read the generated DESIGN.md
  3. Check for "Touch Targets" section
  4. Verify interactive element sizes are listed with compliance status against 44x44pt minimum
- **Expected**: DESIGN.md contains a "Touch Targets" section with element types, their sizes, and yes/no compliance against the 44x44pt guideline.
- **Priority**: P0

## TC-008: Safe area handling inference
- **Source**: Task 2 / AC-4
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/safe-area-handling-inference
- **Pre-conditions**: A target URL is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. Read the generated DESIGN.md
  3. Check for "Safe Areas" section
  4. Verify safe area inset values are listed with source attribution (CSS `env()` / estimated)
- **Expected**: DESIGN.md contains a "Safe Areas" section. If the site uses `env(safe-area-inset-*)` or `viewport-fit=cover`, values are extracted from CSS; otherwise values are marked `(estimated)`.
- **Priority**: P1

## TC-009: Mobile output extends web template with mobile sections
- **Source**: Task 2 / AC-5
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/mobile-output-extends-web-template
- **Pre-conditions**: A target URL is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. Read the generated DESIGN.md
  3. Verify it contains all standard web sections (Color Palette, Typography, Components, etc.)
  4. Verify it also contains mobile-specific sections: Touch Targets, Safe Areas, Responsive Breakpoints
- **Expected**: DESIGN.md has the full web template structure plus three additional mobile-specific sections inserted after "Responsive Behavior" and before "Signature Patterns".
- **Priority**: P0

## TC-010: Mobile match strategy works with closest built-in style
- **Source**: Task 2 / AC-6
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/mobile-match-strategy-built-in
- **Pre-conditions**: A target URL is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <url> --platform mobile`
  2. When prompted for match strategy, select "Match closest built-in style, customize on top"
  3. Verify the output is a mobile-flavored DESIGN.md based on a built-in web style
- **Expected**: Output DESIGN.md is based on the matched built-in style with mobile-specific overrides applied. Contains the "Based on: <style name>" header.
- **Priority**: P1

## TC-011: Mobile output consumable by ui-design
- **Source**: Task 2 / AC-7, Proposal Success Criteria item 4
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/mobile-output-consumable-by-ui-design
- **Pre-conditions**: Mobile DESIGN.md has been generated. The `ui-design` skill is available.
- **Steps**:
  1. Generate DESIGN.md with `extract-design-md <url> --platform mobile`
  2. Invoke `/ui-design` and reference the generated DESIGN.md
  3. Verify ui-design reads and processes the file without errors
- **Expected**: The ui-design skill accepts the mobile DESIGN.md as valid input and generates UI design specifications without modification or manual intervention.
- **Priority**: P0

### TUI Adapter

## TC-012: TUI platform accepts local screenshot path
- **Source**: Task 3 / AC-1
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-accepts-local-screenshot-path
- **Pre-conditions**: A valid terminal screenshot file exists at a local path (e.g. `./testdata/terminal.png`). DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md ./testdata/terminal.png --platform tui`
  2. Verify the command reads the screenshot and proceeds with analysis
  3. Verify DESIGN.md is written
- **Expected**: Command accepts the local file path, reads the screenshot, and generates DESIGN.md without errors.
- **Priority**: P0

## TC-013: TUI rejects URL input with clear error
- **Source**: Task 3 / Hard Rule (TUI input must be local file path)
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-rejects-url-input
- **Pre-conditions**: None.
- **Steps**:
  1. Run `extract-design-md https://example.com --platform tui`
  2. Observe the error output
  3. Verify DESIGN.md was not created
- **Expected**: Command outputs `ERROR: TUI platform requires a local screenshot file path, not a URL. Provide a path like ./screenshot.png` and exits without writing any file.
- **Priority**: P0

## TC-014: TUI AI vision extracts ANSI color palette and character set
- **Source**: Task 3 / AC-2
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-extracts-ansi-colors-and-charset
- **Pre-conditions**: A valid terminal screenshot with visible colors and text is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <screenshot-path> --platform tui`
  2. Read the generated DESIGN.md
  3. Verify "Color Palette" section contains xterm-256 color numbers (0-255 range)
  4. Verify "Character Set" section identifies the character type (box-drawing, block elements, ASCII, or mixed)
  5. Verify "Character Palette Reference" table lists specific characters for borders, dividers, etc.
- **Expected**: DESIGN.md contains ANSI color palette with xterm-256 numbers, character set identification, and a character palette reference table mapping visual elements to specific characters.
- **Priority**: P0

## TC-015: TUI extracts panel layout dimensions and key bindings
- **Source**: Task 3 / AC-2
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-extracts-panel-layout-and-key-bindings
- **Pre-conditions**: A valid terminal screenshot showing multiple panels and a status bar/help legend is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <screenshot-path> --platform tui`
  2. Read the generated DESIGN.md
  3. Verify "Panel Layout" section contains terminal dimensions (rows x columns) and panel dimensions
  4. Verify "Key Bindings" section lists key-to-action mappings if visible in screenshot
- **Expected**: DESIGN.md contains Panel Layout with estimated dimensions and Key Bindings extracted from visible status bar or help panel.
- **Priority**: P0

## TC-016: TUI output matches built-in theme structure
- **Source**: Task 3 / AC-3, Proposal Success Criteria item 5
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-output-matches-builtin-theme-structure
- **Pre-conditions**: A valid terminal screenshot is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <screenshot-path> --platform tui`
  2. Read the generated DESIGN.md
  3. Verify sections match modern-dark-tui structure: Color Space, Character Set, Character Palette Reference, Color Palette, Typography, Panel Layout, Key Bindings, Do's and Don'ts
- **Expected**: DESIGN.md structure aligns with the modern-dark-tui/minimal-ascii-tui template sections, ensuring compatibility with `/ui-design`.
- **Priority**: P0

## TC-017: TUI all values marked as estimated
- **Source**: Task 3 / AC-4, Task 3 / Hard Rule
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-all-values-marked-estimated
- **Pre-conditions**: A valid terminal screenshot is available. DESIGN.md does not exist.
- **Steps**:
  1. Run `extract-design-md <screenshot-path> --platform tui`
  2. Read the generated DESIGN.md
  3. Search for `(estimated)` markers throughout the file
  4. Verify all extracted values (colors, dimensions, characters, key bindings) are marked `(estimated)`
- **Expected**: Every value in the TUI DESIGN.md that was derived from AI vision analysis is annotated with `(estimated)`.
- **Priority**: P0

## TC-018: TUI rejects low-quality screenshot with clear error
- **Source**: Task 3 / AC-6
- **Type**: CLI
- **Target**: cli/extract-design-md
- **Test ID**: cli/extract-design-md/tui-rejects-low-quality-screenshot
- **Pre-conditions**: A blurry or low-resolution image file exists at a local path.
- **Steps**:
  1. Run `extract-design-md <blurry-image-path> --platform tui`
  2. Observe the error output
  3. Verify DESIGN.md was not created
- **Expected**: Command outputs `ERROR: Screenshot quality is too low for reliable analysis. Please provide a clear, high-resolution terminal screenshot.` and exits without writing any file.
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 1 / AC-2, Proposal SC-3 | CLI | cli/extract-design-md | P0 |
| TC-002 | Task 1 / AC-2 | CLI | cli/extract-design-md | P0 |
| TC-003 | Task 1 / AC-4 | CLI | cli/extract-design-md | P0 |
| TC-004 | Task 1 / AC-1 | CLI | cli/extract-design-md | P1 |
| TC-005 | Task 2 / AC-1 | CLI | cli/extract-design-md | P0 |
| TC-006 | Task 2 / AC-2 | CLI | cli/extract-design-md | P0 |
| TC-007 | Task 2 / AC-3 | CLI | cli/extract-design-md | P0 |
| TC-008 | Task 2 / AC-4 | CLI | cli/extract-design-md | P1 |
| TC-009 | Task 2 / AC-5 | CLI | cli/extract-design-md | P0 |
| TC-010 | Task 2 / AC-6 | CLI | cli/extract-design-md | P1 |
| TC-011 | Task 2 / AC-7, Proposal SC-4 | CLI | cli/extract-design-md | P0 |
| TC-012 | Task 3 / AC-1 | CLI | cli/extract-design-md | P0 |
| TC-013 | Task 3 / Hard Rule | CLI | cli/extract-design-md | P0 |
| TC-014 | Task 3 / AC-2 | CLI | cli/extract-design-md | P0 |
| TC-015 | Task 3 / AC-2 | CLI | cli/extract-design-md | P0 |
| TC-016 | Task 3 / AC-3, Proposal SC-5 | CLI | cli/extract-design-md | P0 |
| TC-017 | Task 3 / AC-4, Task 3 / Hard Rule | CLI | cli/extract-design-md | P0 |
| TC-018 | Task 3 / AC-6 | CLI | cli/extract-design-md | P1 |
